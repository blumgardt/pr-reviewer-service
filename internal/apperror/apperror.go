package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeTeamExists  Code = "TEAM_EXISTS"
	CodePRExists    Code = "PR_EXISTS"
	CodePRMerged    Code = "PR_MERGED"
	CodeNotAssigned Code = "NOT_ASSIGNED"
	CodeNoCandidate Code = "NO_CANDIDATE"
	CodeNotFound    Code = "NOT_FOUND"

	CodeValidation Code = "VALIDATION"
	CodeInternal   Code = "INTERNAL"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func From(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
