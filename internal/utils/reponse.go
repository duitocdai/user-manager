package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"           // 400
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"             // 404
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR" // 500
	ErrCodeConfict    ErrorCode = "CONFLICT"              // 409
)

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
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
			"error": appErr.Message,
			"code":  appErr.Code,
		}
		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}
}

func ReponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
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
	default:
		return http.StatusInternalServerError
	}
}
