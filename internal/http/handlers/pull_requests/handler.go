package pull_requests

import (
	"encoding/json"
	"net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/dto/mapping"
	"github.com/blumgardt/pr-reviewer-service.git/internal/http/response"
	"github.com/blumgardt/pr-reviewer-service.git/internal/service"
)

type PullRequestHandler struct {
	prService service.PullRequestService
}

func NewPullRequestHandler(prService service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{prService: prService}
}

// Create godoc
// @Summary      Создать PR и автоматически назначить до 2 ревьюверов
// @Description  Создаёт новый pull request и назначает до двух активных ревьюверов из команды автора.
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreatePullRequestRequest  true  "Pull request create body"
// @Success      201   {object}  dto.CreatePullRequestResponse
// @Failure      400   {object}  response.ErrorResponse             "VALIDATION / NOT_FOUND (author/team)"
// @Failure      409   {object}  response.ErrorResponse             "PR_EXISTS"
// @Router       /pullRequest/create [post]
func (h *PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req dto.CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "invalid json body"))
	}

	pr, err := h.prService.CreatePullRequest(ctx, req.PullRequestID, req.PullRequestID, req.AuthorID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.CreatePullRequestResponse{
		PR: mapping.MapDomainPRToDTO(pr),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// Merge godoc
// @Summary      Пометить PR как MERGED (идемпотентная операция)
// @Description  Переводит pull request в состояние MERGED. Повторный вызов не приводит к ошибке.
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        body  body      dto.MergePullRequestRequest  true  "Pull request id"
// @Success      200   {object}  dto.MergePullRequestResponse
// @Failure      400   {object}  response.ErrorResponse            "VALIDATION"
// @Failure      404   {object}  response.ErrorResponse            "NOT_FOUND"
// @Router       /pullRequest/merge [post]
func (h *PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req dto.MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "invalid json body"))
	}

	mergedPR, err := h.prService.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.MergePullRequestResponse{
		PR: mapping.MapDomainPRToDTO(mergedPR),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// ReAssign godoc
// @Summary      Переназначить ревьювера на другого из его команды
// @Description  Заменяет указанного ревьювера другим активным участником его команды.
// @Tags         PullRequests
// @Accept       json
// @Produce      json
// @Param        body  body      dto.ReassignPullRequestRequest  true  "Reassign body"
// @Success      200   {object}  dto.ReassignPullRequestResponse
// @Failure      400   {object}  response.ErrorResponse               "VALIDATION"
// @Failure      404   {object}  response.ErrorResponse               "NOT_FOUND"
// @Failure      409   {object}  response.ErrorResponse               "PR_MERGED / NOT_ASSIGNED / NO_CANDIDATE"
// @Router       /pullRequest/reassign [post]
func (h *PullRequestHandler) ReAssign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req dto.ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, apperror.New(apperror.CodeValidation, "invalid json body"))
	}

	reAssignedPR, replacedBy, err := h.prService.ReAssignPullRequest(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	resp := dto.ReassignPullRequestResponse{
		PR:         mapping.MapDomainPRToDTO(reAssignedPR),
		ReplacedBy: replacedBy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
