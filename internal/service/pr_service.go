package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/repository/postgres"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error)
	MergePullRequest(ctx context.Context, id string) (*domain.PullRequest, error)
	ReAssignPullRequest(ctx context.Context, id, oldUserID string) (*domain.PullRequest, string, error)
}

type pullRequestService struct {
	prRepo   postgres.PullRequestRepository
	userRepo postgres.UserRepository
	teamRepo postgres.TeamRepository
}

func NewPullRequestService(prRepo postgres.PullRequestRepository, userRepo postgres.UserRepository, teamRepo postgres.TeamRepository) PullRequestService {
	return &pullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *pullRequestService) CreatePullRequest(
	ctx context.Context,
	id, name, authorID string,
) (*domain.PullRequest, error) {
	if id == "" || name == "" || authorID == "" {
		return nil, apperror.New(apperror.CodeValidation, "pull_request_id, pull_request_name and author_id are required")
	}

	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if author.TeamName == "" {
		return nil, apperror.New(apperror.CodeValidation, "author has no teams")
	}
	teamName := author.TeamName

	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	reviewersID := s.pickReviewers(team, authorID)

	pr := &domain.PullRequest{
		PullRequestID:     id,
		PullRequestName:   name,
		AuthorID:          authorID,
		PullRequestStatus: string(domain.PRStatusOpen),
		ReviewersID:       reviewersID,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *pullRequestService) MergePullRequest(ctx context.Context, id string) (*domain.PullRequest, error) {
	if id == "" {
		return nil, apperror.New(apperror.CodeValidation, "pull_request_id is required")
	}

	pr, err := s.prRepo.Merge(ctx, id)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *pullRequestService) ReAssignPullRequest(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
	if prID == "" || oldUserID == "" {
		return nil, "", apperror.New(apperror.CodeValidation, "pull_request_id and old_user_id are required")
	}

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}

	if pr.PullRequestStatus == string(domain.PRStatusMerged) {
		return nil, "", apperror.New(apperror.CodePRMerged, "cannot reassign on merged PR")
	}

	isAssigned := false
	for _, rid := range pr.ReviewersID {
		if rid == oldUserID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", apperror.New(apperror.CodeNotAssigned, "reviewer is not assigned to this PR")
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return nil, "", err
	}
	if oldReviewer.TeamName == "" {
		return nil, "", apperror.New(apperror.CodeValidation, "old reviewer has no teams")
	}

	team, err := s.teamRepo.GetTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	assignedSet := make(map[string]struct{})
	for _, rid := range pr.ReviewersID {
		assignedSet[rid] = struct{}{}
	}

	var candidates []string
	for _, m := range team.Members {
		if !m.IsActive {
			continue
		}
		if m.ID == oldUserID {
			continue
		}
		if m.ID == pr.AuthorID {
			continue
		}
		if _, alreadyAssigned := assignedSet[m.ID]; alreadyAssigned {
			continue
		}
		candidates = append(candidates, m.ID)
	}

	if len(candidates) == 0 {
		return nil, "", apperror.New(apperror.CodeNoCandidate, "no active replacement candidate in teams")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newIdx := r.Intn(len(candidates))
	newReviewerID := candidates[newIdx]

	updatedPR, err := s.prRepo.ReAssign(ctx, prID, oldUserID, newReviewerID)
	if err != nil {
		return nil, "", err
	}

	return updatedPR, newReviewerID, nil
}

func (s *pullRequestService) pickReviewers(team *domain.Team, authorID string) []string {
	var candidates []string
	for _, m := range team.Members {
		if !m.IsActive {
			continue
		}
		if m.ID == authorID {
			continue
		}
		candidates = append(candidates, m.ID)
	}

	if len(candidates) <= 2 {
		out := make([]string, len(candidates))
		copy(out, candidates)
		return out
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:2]
}
