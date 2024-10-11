package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type EchoError struct {
	Code     int       `json:"code"`
	Message  string    `json:"message"`
	Error    string    `json:"error"`
	IssuedAt time.Time `json:"issued_at"`
}

type ValidationErrors struct {
	Code     int       `json:"code"`
	Message  string    `json:"message"`
	Error    string    `json:"error"`
	IssuedAt time.Time `json:"issued_at"`
	Errors   []string  `json:"errors"`
}

func NewInternalServerError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusInternalServerError,
		Message:  "Internal Server Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

func NewConflictError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusConflict,
		Message:  "Conflict Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

func NewNotFoundError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusNotFound,
		Message:  "Not Found Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

func NewUnauthorizedError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusUnauthorized,
		Message:  "Unauthorized Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

func NewBadRequestError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusBadRequest,
		Message:  "Bad Request Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

func NewUnprocessableEntityError(message string) error {
	jsonError := EchoError{
		Code:     http.StatusUnprocessableEntity,
		Message:  "Unprocessable Entity Error",
		Error:    message,
		IssuedAt: time.Now(),
	}
	return echo.NewHTTPError(jsonError.Code, jsonError)
}

// Show the log and return true if there was an error
func FailOnError(logger *logrus.Entry, err error) bool {
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return true
	}
	return false
}

func WarnOnError(logger *logrus.Entry, err error) bool {
	if err != nil {
		logger.WithError(err).Warn(err.Error())
		return true
	}
	return false
}

func DebugOnError(logger *logrus.Entry, err error) bool {
	if err != nil {
		logger.WithError(err).Debug(err.Error())
		return true
	}
	return false
}
