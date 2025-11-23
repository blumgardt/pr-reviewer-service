package service

import (
	"context"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/repository/postgres"
)

type TeamService interface {
	Add(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}

type teamService struct {
	teamRepo postgres.TeamRepository
}

func NewTeamService(teamRepo postgres.TeamRepository) TeamService {
	return &teamService{teamRepo: teamRepo}
}

func (s *teamService) Add(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	if team == nil {
		return nil, apperror.New(apperror.CodeValidation, "teams is required")
	}
	if team.Name == "" {
		return nil, apperror.New(apperror.CodeValidation, "team_name is required")
	}
	if len(team.Members) == 0 {
		return nil, apperror.New(apperror.CodeValidation, "teams must contain at least one member")
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *teamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	if name == "" {
		return nil, apperror.New(apperror.CodeValidation, "team_name is required")
	}

	team, err := s.teamRepo.GetTeam(ctx, name)
	if err != nil {
		return nil, err
	}

	return team, nil
}
