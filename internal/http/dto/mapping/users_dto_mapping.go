package mapping

import (
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
)

func MapDomainUserToDTO(u *domain.User) dto.UserDTO {
	return dto.UserDTO{
		UserID:   u.ID,
		Username: u.Name,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
