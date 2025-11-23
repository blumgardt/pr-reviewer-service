package stats

import (
	"encoding/json"
	"net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto/mapping"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/response"
	"github.com/blumgardt/pr-reviewer-service.git/internal/service"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// GetReviewerStats godoc
// @Summary      Статистика по ревьюверам
// @Description  Возвращает количество назначений на ревью для каждого пользователя
// @Tags         Stats
// @Produce      json
// @Success      200  {object}  dto.ReviewerStatsResponse
// @Failure      500  {object}  response.ErrorResponse   "INTERNAL"
// @Router       /stats/reviewers [get]
func (h *StatsHandler) GetReviewerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.statsService.GetReviewerStats(ctx)
	if err != nil {
		response.WriteError(w, apperror.From(err))
		return
	}

	resp := mapping.ReviewerStatsToDTO(stats)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
