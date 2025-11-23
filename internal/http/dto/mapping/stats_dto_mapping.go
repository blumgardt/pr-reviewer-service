package mapping

import (
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
)

func ReviewerStatsToDTO(stats []domain.ReviewerStats) dto.ReviewerStatsResponse {
	resp := dto.ReviewerStatsResponse{
		Items: make([]dto.ReviewerStatsItem, 0, len(stats)),
	}
	for _, s := range stats {
		resp.Items = append(resp.Items, dto.ReviewerStatsItem{
			UserID:        s.UserID,
			Username:      s.UserName,
			AssignedCount: s.AssignedCount,
		})
	}
	return resp
}
