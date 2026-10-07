// Package services 实现核心业务逻辑。
package services

import (
	"errors"
	"net/http"
)

// AppError 是带 HTTP 状态码的业务错误。
type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string { return e.Message }

// NewAppError 构造业务错误。
func NewAppError(status int, msg string) *AppError {
	return &AppError{Status: status, Message: msg}
}

// ErrNotFound 构造 404。
func ErrNotFound(msg string) *AppError { return NewAppError(http.StatusNotFound, msg) }

// ErrBadRequest 构造 400。
func ErrBadRequest(msg string) *AppError { return NewAppError(http.StatusBadRequest, msg) }

// ErrUnprocessable 构造 422（业务守卫）。
func ErrUnprocessable(msg string) *AppError { return NewAppError(http.StatusUnprocessableEntity, msg) }

// ErrForbidden 构造 403（已登录但无权限操作该资源）。
func ErrForbidden(msg string) *AppError { return NewAppError(http.StatusForbidden, msg) }

// AsAppError 提取 AppError。
func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// ValidationError 带字段错误的 422（Laravel validator 风格：message="参数错误" + errors 键值表）。
// 由 handler 的 fail() 统一渲染为 `{"message":"参数错误","errors":{...}}`。
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string { return "参数错误" }

// NewValidationError 构造带字段错误的 422。
func NewValidationError(errs map[string][]string) *ValidationError {
	return &ValidationError{Errors: errs}
}

// AsValidationError 提取 ValidationError。
func AsValidationError(err error) (*ValidationError, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ve, true
	}
	return nil, false
}
