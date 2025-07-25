package errors

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code         int    `json:"code"`    // Включить код ошибки в JSON
	Message      string `json:"message"` // Сообщение для клиента
	Err          error  `json:"-"`       // Не отправлять внутренние ошибки
	IsUserFacing bool   `json:"-"`       // Внутреннее поле
}

func (a *AppError) Error() string {
	if a == nil {
		return ""
	}
	if a.Err != nil {
		return fmt.Sprintf("%s (code: %d): %v", a.Message, a.Code, a.Err)
	}
	return fmt.Sprintf("%s (code: %d)", a.Message, a.Code)
}

type DBError struct {
	Message string
	Err     error
}

const (
	InternalServerError = "internal server error"
	BadRequest          = "bad request"
	NotFound            = "not_found"

	InvalidDataCode         = 402
	InternalServerErrorCode = 500
	NotFoundErrorCode       = 404
)

func NewAppError(httpCode int, message string, err error, isUserFacing bool) *AppError {
	return &AppError{
		Code:         httpCode,
		Message:      message,
		Err:          err,
		IsUserFacing: isUserFacing,
	}
}

func NewDBError(message string, dbError error) *AppError {
	return &AppError{
		Code:         InternalServerErrorCode,
		Message:      message,
		Err:          dbError,
		IsUserFacing: false,
	}
}

var (
	ErrEmptyAction  = errors.New("action did not affect the data")
	ErrDataNotFound = errors.New("data not found")
	ErrEmptyData    = errors.New("empty data")
)

func Is(err any, err2 error) bool {
	if e, ok := err.(error); ok {
		return errors.Is(e, err2)
	}
	return false
}

var ErrNotFound = errors.New("not found")

func NewNotFoundError(message string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, message)
}
