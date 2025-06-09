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
	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (e *ErrorHandler) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Warnw("bad request", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (e *ErrorHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	e.logger.Warnw("not found", "method", r.Method, "url", r.URL.Path)
	WriteJSONError(w, http.StatusNotFound, "not found")
}

func (e *ErrorHandler) Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	e.logger.Warnw("unauthorized", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}
