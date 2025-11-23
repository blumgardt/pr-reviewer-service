package mapping

import (
	"github.com/blumgardt/pr-reviewer-service.git/internal/domain"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
)

func MapTeamDTOToDomain(dto *dto.TeamDTO) *domain.Team {
	team := &domain.Team{
		Name: dto.TeamName,
	}

	for _, m := range dto.Members {
		team.Members = append(team.Members, domain.User{
			ID:       m.UserID,
			Name:     m.Username,
			IsActive: m.IsActive,
			TeamName: dto.TeamName,
		})
	}

	return team
}

func MapDomainTeamToDTO(team *domain.Team) dto.TeamDTO {
	newDTO := dto.TeamDTO{
		TeamName: team.Name,
	}

	for _, m := range team.Members {
		newDTO.Members = append(newDTO.Members, dto.TeamMemberDTO{
			UserID:   m.ID,
			Username: m.Name,
			IsActive: m.IsActive,
		})
	}

	return newDTO
}
