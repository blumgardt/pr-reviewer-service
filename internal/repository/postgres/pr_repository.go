package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository interface {
	Create(ctx context.Context, request *domain.PullRequest) error
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	ReAssign(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*domain.PullRequest, error)
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
}

type pullRequestRepository struct {
	db *pgxpool.Pool
}

func NewPullRequestRepository(dbPool *pgxpool.Pool) PullRequestRepository {
	return &pullRequestRepository{db: dbPool}
}

func (r *pullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) (err error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "begin tx", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const insertPR = `
		INSERT INTO pull_requests (id, name, author_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, merged_at;
	`

	if err = tx.QueryRow(ctx, insertPR,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.PullRequestStatus,
	).Scan(&pr.CreatedAt, &pr.MergedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperror.New(apperror.CodePRExists, "pull_request_id already exists")
		}
		return apperror.Wrap(apperror.CodeInternal, "insert pull_request", err)
	}

	if len(pr.ReviewersID) > 0 {
		const insertReviewer = `
			INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
			VALUES ($1, $2);
		`
		for _, rid := range pr.ReviewersID {
			if _, err = tx.Exec(ctx, insertReviewer, pr.PullRequestID, rid); err != nil {
				return apperror.Wrap(apperror.CodeInternal, fmt.Sprintf("insert reviewer %s", rid), err)
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "commit tx", err)
	}

	return nil
}

func (r *pullRequestRepository) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	const q = `
		UPDATE pull_requests
		SET status   = 'MERGED',
		    merged_at = COALESCE(merged_at, now())
		WHERE id = $1
		RETURNING id, name, author_id, status, created_at, merged_at;
	`

	var pr domain.PullRequest

	err := r.db.QueryRow(ctx, q, id).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.PullRequestStatus,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(apperror.CodeNotFound, "pull request not found")
		}
		return nil, apperror.Wrap(apperror.CodeInternal, "merge pull request", err)
	}

	const reviewersQuery = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
		ORDER BY reviewer_id;
	`

	rows, err := r.db.Query(ctx, reviewersQuery, id)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "query reviewers", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "scan reviewer_id", err)
		}
		pr.ReviewersID = append(pr.ReviewersID, rid)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "iterate reviewers", err)
	}

	return &pr, nil
}

func (r *pullRequestRepository) ReAssign(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*domain.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "begin tx", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const updateReviewer = `
        UPDATE pull_request_reviewers
        SET reviewer_id = $3
        WHERE pull_request_id = $1 AND reviewer_id = $2;
    `

	tag, execErr := tx.Exec(ctx, updateReviewer, prID, oldReviewerID, newReviewerID)
	if execErr != nil {
		err = apperror.Wrap(apperror.CodeInternal, "update reviewer", execErr)
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		err = apperror.New(apperror.CodeNotAssigned, "reviewer is not assigned to this PR")
		return nil, err
	}

	const selectPR = `
        SELECT id, name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE id = $1;
    `

	var pr domain.PullRequest
	err = tx.QueryRow(ctx, selectPR, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.PullRequestStatus,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		err = apperror.Wrap(apperror.CodeInternal, "select pull request", err)
		return nil, err
	}

	const reviewersQuery = `
        SELECT reviewer_id
        FROM pull_request_reviewers
        WHERE pull_request_id = $1
        ORDER BY reviewer_id;
    `
	rows, err := tx.Query(ctx, reviewersQuery, prID)
	if err != nil {
		err = apperror.Wrap(apperror.CodeInternal, "query reviewers", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rid string
		if scanErr := rows.Scan(&rid); scanErr != nil {
			err = apperror.Wrap(apperror.CodeInternal, "scan reviewer_id", scanErr)
			return nil, err
		}
		pr.ReviewersID = append(pr.ReviewersID, rid)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		err = apperror.Wrap(apperror.CodeInternal, "iterate reviewers", rowsErr)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		err = apperror.Wrap(apperror.CodeInternal, "commit tx", err)
		return nil, err
	}

	return &pr, nil
}

func (r *pullRequestRepository) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	const q = `
        SELECT id, name, author_id, status, created_at, merged_at
        FROM pull_requests
        WHERE id = $1;
    `

	var pr domain.PullRequest

	err := r.db.QueryRow(ctx, q, id).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.PullRequestStatus,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(apperror.CodeNotFound, "pull request not found")
		}
		return nil, apperror.Wrap(apperror.CodeInternal, "get pull request", err)
	}

	const reviewersQuery = `
        SELECT reviewer_id
        FROM pull_request_reviewers
        WHERE pull_request_id = $1
        ORDER BY reviewer_id;
    `
	rows, err := r.db.Query(ctx, reviewersQuery, id)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "query reviewers", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rid string
		if err := rows.Scan(&rid); err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "scan reviewer_id", err)
		}
		pr.ReviewersID = append(pr.ReviewersID, rid)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "iterate reviewers", err)
	}

	return &pr, nil
}
