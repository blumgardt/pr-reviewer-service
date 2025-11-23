package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, id, name string, isActive bool) (*domain.User, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateActiveStatus(ctx context.Context, user *domain.User, isActive bool) error
	GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{db: pool}
}

func (r *userRepository) CreateNewUser(ctx context.Context, id, name string, isActive bool) (*domain.User, error) {
	const q = `
		INSERT INTO user.users (id, name, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, name, is_active
    `

	var usr domain.User

	err := r.db.QueryRow(ctx, q, id, name, isActive).Scan(
		&usr.ID,
		&usr.Name,
		&usr.IsActive,
	)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "create user", err)
	}

	return &usr, nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	const query = `
        SELECT id, name, is_active, team_name
        FROM users
        WHERE id = $1
    `

	var u domain.User
	var teamName string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID,
		&u.Name,
		&u.IsActive,
		&teamName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.New(apperror.CodeNotFound, "user not found")
		}
		return nil, apperror.Wrap(apperror.CodeInternal, "get user by id", err)
	}

	u.TeamName = teamName
	return &u, nil
}

func (r *userRepository) UpdateActiveStatus(ctx context.Context, user *domain.User, isActive bool) error {
	const q = `
		UPDATE users
		SET is_active = $2
		WHERE id = $1
    `

	_, err := r.db.Exec(ctx, q, user.ID, isActive)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "update user active status", err)
	}

	return nil
}

func (r *userRepository) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	const q = `
		SELECT pr.id, pr.name, pr.author_id, pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN pull_request_reviewers r
		  ON r.pull_request_id = pr.id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at DESC;
	`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "query pull requests by reviewer", err)
	}
	defer rows.Close()

	var result []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		var mergedAt *time.Time

		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.PullRequestStatus,
			&pr.CreatedAt,
			&mergedAt,
		); err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "scan pull request", err)
		}

		pr.MergedAt = mergedAt
		result = append(result, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "iterate pull requests", err)
	}

	return result, nil
}
