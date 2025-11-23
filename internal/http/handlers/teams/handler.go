package teams

import (
	"encoding/json"
	"net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto/mapping"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/response"
	"github.com/blumgardt/pr-reviewer-service.git/internal/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// Add godoc
// @Summary      Создать команду с участниками (создаёт/обновляет пользователей)
// @Description  Создаёт новую команду с участниками. Если пользователи уже существуют, обновляет их данные и привязку к команде.
// @Tags         Teams
// @Accept       json
// @Produce      json
// @Param        body  body      dto.TeamDTO           true  "Team body"
// @Success      201   {object}  dto.TeamAddResponse
// @Failure      400   {object}  response.ErrorResponse     "VALIDATION"
// @Failure      409   {object}  response.ErrorResponse     "TEAM_EXISTS"
// @Router       /team/add [post]
func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req dto.TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "invalid JSON body"))
		return
	}

	team := mapping.MapTeamDTOToDomain(&req)

	created, err := h.teamService.Add(ctx, team)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.TeamAddResponse{
		Team: mapping.MapDomainTeamToDTO(created),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// Get godoc
// @Summary      Получить команду с участниками
// @Description  Возвращает команду и список её участников по имени команды.
// @Tags         Teams
// @Produce      json
// @Param        team_name  query     string        true  "Team name"
// @Success      200        {object}  dto.TeamDTO
// @Failure      400        {object}  response.ErrorResponse   "VALIDATION"
// @Failure      404        {object}  response.ErrorResponse   "NOT_FOUND"
// @Router       /team/get [get]
func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "team_name is required"))
		return
	}

	team, err := h.teamService.GetTeam(ctx, teamName)
	if err != nil {
		response.WriteError(w, err)
	}

	resp := mapping.MapDomainTeamToDTO(team)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
