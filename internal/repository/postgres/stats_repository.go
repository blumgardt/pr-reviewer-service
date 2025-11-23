package postgres

import (
	"context"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository interface {
	GetReviewerStats(ctx context.Context) ([]domain.ReviewerStats, error)
}

type statsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) GetReviewerStats(ctx context.Context) ([]domain.ReviewerStats, error) {
	const q = `
		SELECT u.id, u.name, COUNT(prr.pull_request_id) AS assigned_count
		FROM users u
		LEFT JOIN pull_request_reviewers prr
		  ON prr.reviewer_id = u.id
		GROUP BY u.id, u.name
		ORDER BY assigned_count DESC, u.name;
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "query reviewer stats", err)
	}
	defer rows.Close()

	var res []domain.ReviewerStats

	for rows.Next() {
		var s domain.ReviewerStats
		if err := rows.Scan(&s.UserID, &s.UserName, &s.AssignedCount); err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "scan reviewer stats", err)
		}
		res = append(res, s)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "iterate reviewer stats", err)
	}

	return res, nil
}
