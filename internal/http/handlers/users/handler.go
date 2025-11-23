package users

import (
	"encoding/json"
	"net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto/mapping"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/response"
	"github.com/blumgardt/pr-reviewer-service.git/internal/service"
)

type UsersHandler struct {
	userService service.UserService
}

func NewUsersHandler(userService service.UserService) *UsersHandler {
	return &UsersHandler{userService: userService}
}

func (h *UsersHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req dto.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "invalid json body"))
		return
	}

	user, err := h.userService.SetIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.SetIsActiveResponse{
		User: mapping.MapDomainUserToDTO(user),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *UsersHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")

	prs, err := h.userService.GetReview(ctx, userID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.GetReviewResponse{
		UserID:       userID,
		PullRequests: make([]dto.PullRequestShortDTO, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, mapping.MapDomainPRToShortDTO(pr))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
