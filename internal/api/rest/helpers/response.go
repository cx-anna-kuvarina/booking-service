package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorCode int

const (
	InvalidQueries ErrorCode = iota
	InvalidStateParameter
	ExchangeTokenErr
	AuthUserInfoErr
	MissingAuthHeaderErr
	InvalidHeaderFormatErr
	InvalidTokenErr
	DecodeUserInfoErr
	GenerateTokenErr
	GetUserIdErr
	InvalidEmailErr
	InvalidUserInfoErr
	UsersStoreErr
	NotFoundErr
	ReadBodyErr
	InvalidRequest
	ValidationError
	NotFound
	InternalError
)

type ErrorResponse struct {
	Message string
	Type    string
	Code    ErrorCode
}

// ValidationErr represents a validation error
type ValidationErr struct {
	Message string
}

func (e *ValidationErr) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationErr{Message: message}
}

func WriteErrorResponse(w http.ResponseWriter, resp *ErrorResponse, status int) {
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("failed to parse error response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(jsonStr)
}

func NewErrorResponse(message string, code ErrorCode) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Type:    "ERROR",
		Code:    code,
	}
}

func WriteData(ctx context.Context, w http.ResponseWriter, response interface{}, httpStatus int) {
	var toMarshal any = response
	var err error

	respJSON := make([]byte, 0)
	if toMarshal != nil {
		respJSON, err = json.Marshal(toMarshal)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed marshaling JSON response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
	}

	w.WriteHeader(httpStatus)
	_, _ = w.Write(respJSON)
}
