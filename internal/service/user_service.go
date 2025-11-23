package service

import (
	"context"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/repository/postgres"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type userService struct {
	userRepo postgres.UserRepository
}

func NewUserService(userRepo postgres.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	if userID == "" {
		return nil, apperror.New(apperror.CodeValidation, "user_id is required")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateActiveStatus(ctx, user, isActive); err != nil {
		return nil, err
	}

	user.IsActive = isActive
	return user, nil
}

func (s *userService) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	if userID == "" {
		return nil, apperror.New(apperror.CodeValidation, "user_id is required")
	}

	pullRequests, err := s.userRepo.GetReview(ctx, userID)
	if err != nil {
		return nil, err
	}

	return pullRequests, nil
}
