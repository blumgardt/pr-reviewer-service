package response

import (
	"encoding/json"
	"net/http"

	"github.com/blumgardt/pr-reviewer-service.git/internal/apperror"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func WriteError(w http.ResponseWriter, err error) {
	appErr := apperror.From(err)

	status := http.StatusInternalServerError
	code := "INTERNAL"
	message := "internal error"

	if appErr != nil {
		code = string(appErr.Code)
		message = appErr.Message

		switch appErr.Code {
		case apperror.CodeValidation:
			status = http.StatusBadRequest
		case apperror.CodeNotFound:
			status = http.StatusNotFound
		case apperror.CodeTeamExists,
			apperror.CodePRExists,
			apperror.CodePRMerged,
			apperror.CodeNotAssigned,
			apperror.CodeNoCandidate:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message

	_ = json.NewEncoder(w).Encode(resp)
}
