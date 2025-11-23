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

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}

type teamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) (err error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "begin tx", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const insertTeam = `INSERT INTO teams (name) VALUES ($1)`

	if _, err = tx.Exec(ctx, insertTeam, team.Name); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return apperror.New(apperror.CodeTeamExists, "team_name already exists")
		}
		return apperror.Wrap(apperror.CodeInternal, "insert teams", err)
	}

	const upsertUser = `
		INSERT INTO users (id, name, is_active, team_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = excluded.name,
			is_active = excluded.is_active,
			team_name = excluded.team_name;
	`

	for _, m := range team.Members {
		_, err = tx.Exec(ctx, upsertUser,
			m.ID,
			m.Name,
			m.IsActive,
			team.Name,
		)
		if err != nil {
			return apperror.Wrap(apperror.CodeInternal, fmt.Sprintf("upsert user %s", m.ID), err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "commit tx", err)
	}

	return nil
}

func (r *teamRepository) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	const q = `
    	SELECT id, name, is_active, team_name
		FROM users
		WHERE team_name = $1
    `

	rows, err := r.db.Query(ctx, q, name)
	if err != nil {
		return nil, fmt.Errorf("query teams %s: $w", name, err)
	}
	defer rows.Close()

	var result domain.Team

	for rows.Next() {
		var member domain.User

		if err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.IsActive,
			&member.TeamName,
		); err != nil {
			return nil, fmt.Errorf("scan teams member: %w", err)
		}

		result.Members = append(result.Members, member)
	}

	result.Name = result.Members[0].TeamName

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams members: %w", err)
	}

	return &result, nil
}
