package service

import (
	"context"

	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/repository/postgres"
)

type StatsService interface {
	GetReviewerStats(ctx context.Context) ([]domain.ReviewerStats, error)
}

type statsService struct {
	statsRepo postgres.StatsRepository
}

func NewStatsService(statsRepo postgres.StatsRepository) StatsService {
	return &statsService{statsRepo: statsRepo}
}

func (s *statsService) GetReviewerStats(ctx context.Context) ([]domain.ReviewerStats, error) {
	return s.statsRepo.GetReviewerStats(ctx)
}
