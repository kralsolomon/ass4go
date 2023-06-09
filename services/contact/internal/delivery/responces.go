package delivery

import (
	"encoding/json"
	"net/http"

	"advanced.microservices/pkg/jsonlog"
)

// func logError(r *http.Request, err error) {
// 	logger.PrintError(err, map[string]string{
// 		"request_method": r.Method,
// 		"request_url":    r.URL.String(),
// 	})
// }

type envelope map[string]any
type responseHandler struct {
	logger *jsonlog.Logger
}

func writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (handler *responseHandler) logError(r *http.Request, err error) {
	handler.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (handler *responseHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := writeJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(status)
	}
}

func (handler *responseHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	handler.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	handler.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (handler *responseHandler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	handler.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (handler *responseHandler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	handler.errorResponse(w, r, http.StatusNotFound, message)
}

func (handler *responseHandler) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	handler.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (handler *responseHandler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	handler.errorResponse(w, r, http.StatusConflict, message)
}
