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
