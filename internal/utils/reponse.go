package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"           // 400
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"             // 404
	ErrCodeInternal        ErrorCode = "INTERNAL_SERVER_ERROR" // 500
	ErrCodeConfict         ErrorCode = "CONFLICT"              // 409
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrrCodeTooManyRequest ErrorCode = "TOO_MANY_REQUEST"
)

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}
type APIResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
}

func (ae *AppError) Error() string {
	return ae.Message
}
func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}
func WrapError(err error, message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func ReponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": CapitalizeFirst(appErr.Message),
			"code":  appErr.Code,
		}
		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}
}

func ReponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	resp := APIResponse{
		Status:  "success",
		Message: CapitalizeFirst(message),
	}

	if len(data) > 0 && data[0] != nil {
		if m, ok := data[0].(map[string]any); ok {
			if p, exists := m["pagination"]; exists {
				resp.Pagination = p
			}

			if d, exists := m["data"]; exists {
				resp.Data = d
			} else {
				resp.Data = m
			}
		} else {
			resp.Data = data[0]
		}
	}
	ctx.JSON(status, resp)
}
func ReponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

func ReponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeInternal:
		return http.StatusInternalServerError
	case ErrCodeConfict:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrrCodeTooManyRequest:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
