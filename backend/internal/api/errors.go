package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Validation errors
	ErrCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrCodeInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrCodeMissingParameter ErrorCode = "MISSING_PARAMETER"

	// Resource errors
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"

	// Content specific errors
	ErrCodeContentNotFound ErrorCode = "CONTENT_NOT_FOUND"
	ErrCodeProviderError   ErrorCode = "PROVIDER_ERROR"

	// System errors
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabaseError  ErrorCode = "DATABASE_ERROR"
	ErrCodeCacheError     ErrorCode = "CACHE_ERROR"
	ErrCodeRateLimitError ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Admin errors
	ErrCodeAdminRequired ErrorCode = "ADMIN_ACCESS_REQUIRED"
	ErrCodeInvalidAPIKey ErrorCode = "INVALID_API_KEY"
)

// Error represents a structured API error
type Error struct {
	Code      ErrorCode         `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Timestamp string            `json:"timestamp,omitempty"`
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Message
}

// NewError creates a new structured error
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: make(map[string]string),
	}
}

// WithDetails adds details to the error
func (e *Error) WithDetails(key, value string) *Error {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// WithRequestID adds request ID to the error
func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

// WithTimestamp adds timestamp to the error
func (e *Error) WithTimestamp(timestamp string) *Error {
	e.Timestamp = timestamp
	return e
}

// ErrorResponse represents the standard API error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err *Error) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
	}
}

// ErrorHandler is a Gin middleware for handling structured errors
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If there are no errors, continue
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last()

		// Try to cast to our structured error
		if apiErr, ok := err.Err.(*Error); ok {
			// It's already a structured error
			statusCode := getHTTPStatusForErrorCode(apiErr.Code)
			c.JSON(statusCode, NewErrorResponse(apiErr))
			return
		}

		// It's a regular error, wrap it
		logger.Error("unhandled error",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err))

		apiErr := NewError(ErrCodeInternalError, "An unexpected error occurred").
			WithDetails("original_error", err.Error())

		c.JSON(http.StatusInternalServerError, NewErrorResponse(apiErr))
	}
}

// SendError sends a structured error response
func SendError(c *gin.Context, err *Error) {
	statusCode := getHTTPStatusForErrorCode(err.Code)
	c.AbortWithStatusJSON(statusCode, NewErrorResponse(err))
}

// getHTTPStatusForErrorCode maps error codes to HTTP status codes
func getHTTPStatusForErrorCode(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeInvalidParameter, ErrCodeMissingParameter:
		return http.StatusBadRequest
	case ErrCodeNotFound, ErrCodeContentNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeForbidden, ErrCodeAdminRequired:
		return http.StatusForbidden
	case ErrCodeUnauthorized, ErrCodeInvalidAPIKey:
		return http.StatusUnauthorized
	case ErrCodeRateLimitError:
		return http.StatusTooManyRequests
	case ErrCodeInternalError, ErrCodeDatabaseError, ErrCodeCacheError, ErrCodeProviderError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ErrInvalidParameter creates an error for invalid parameters
func ErrInvalidParameter(param, reason string) *Error {
	return NewError(ErrCodeInvalidParameter, "Invalid parameter").
		WithDetails("parameter", param).
		WithDetails("reason", reason)
}

func ErrMissingParameter(param string) *Error {
	return NewError(ErrCodeMissingParameter, "Missing required parameter").
		WithDetails("parameter", param)
}

func ErrContentNotFound(id string) *Error {
	return NewError(ErrCodeContentNotFound, "Content not found").
		WithDetails("content_id", id)
}

func ErrUnauthorized() *Error {
	return NewError(ErrCodeUnauthorized, "Authentication required")
}

func ErrForbidden() *Error {
	return NewError(ErrCodeForbidden, "Access forbidden")
}

func ErrInternal(message string) *Error {
	return NewError(ErrCodeInternalError, message)
}

func ErrRateLimitExceeded() *Error {
	return NewError(ErrCodeRateLimitError, "Rate limit exceeded")
}

func ErrInvalidAPIKey() *Error {
	return NewError(ErrCodeInvalidAPIKey, "Invalid API key")
}
