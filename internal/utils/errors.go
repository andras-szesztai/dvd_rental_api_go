package utils

import (
	"net/http"

	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.SugaredLogger
}

func NewErrorHandler(logger *zap.SugaredLogger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (e *ErrorHandler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Errorw("internal server error", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	err = WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
	if err != nil {
		e.logger.Errorw("failed to write JSON error", "error", err.Error())
	}
}

func (e *ErrorHandler) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Warnw("bad request", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	err = WriteJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		e.logger.Errorw("failed to write JSON error", "error", err.Error())
	}
}

func (e *ErrorHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	e.logger.Warnw("not found", "method", r.Method, "url", r.URL.Path)
	err := WriteJSONError(w, http.StatusNotFound, "not found")
	if err != nil {
		e.logger.Errorw("failed to write JSON error", "error", err.Error())
	}
}

func (e *ErrorHandler) Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Warnw("unauthorized", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	err = WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
	if err != nil {
		e.logger.Errorw("failed to write JSON error", "error", err.Error())
	}
}

func (e *ErrorHandler) TooManyRequests(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Warnw("too many requests", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	err = WriteJSONError(w, http.StatusTooManyRequests, "too many requests")
	if err != nil {
		e.logger.Errorw("failed to write JSON error", "error", err.Error())
	}
}
