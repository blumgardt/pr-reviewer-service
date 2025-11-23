package mapping

import (
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
)

func MapDomainPRToDTO(pr *domain.PullRequest) dto.PullRequestDTO {
	return dto.PullRequestDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.PullRequestStatus,
		AssignedReviewers: pr.ReviewersID,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func MapDomainPRToShortDTO(pr domain.PullRequest) dto.PullRequestShortDTO {
	return dto.PullRequestShortDTO{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.PullRequestStatus,
	}
}
