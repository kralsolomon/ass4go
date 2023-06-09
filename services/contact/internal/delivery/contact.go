package delivery

import (
	"errors"
	"fmt"
	"net/http"

	"advanced.microservices/pkg/helpers"
	"advanced.microservices/pkg/jsonlog"
	"advanced.microservices/pkg/validator"
	"advanced.microservices/services/contact/internal/domain"
	"advanced.microservices/services/contact/internal/repository"
	"github.com/julienschmidt/httprouter"
)

type ContactHandler struct {
	contactUseCase domain.ContactUseCase
	response       responseHandler
}

func NewContactHandler(router *httprouter.Router, logger *jsonlog.Logger, contactUseCase domain.ContactUseCase) {
	handler := &ContactHandler{
		contactUseCase: contactUseCase,
		response:       responseHandler{logger: logger},
	}
	router.HandlerFunc(http.MethodGet, "/contact", handler.getById)
	router.HandlerFunc(http.MethodPost, "/contact", handler.create)
	router.HandlerFunc(http.MethodDelete, "/contact", handler.delete)
	router.HandlerFunc(http.MethodPut, "/contact", handler.update)
	router.HandlerFunc(http.MethodGet, "/contact/healthcheck", handler.healthcheck)
}

func (handler *ContactHandler) healthcheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, envelope{"status": "ok"}, nil)
}

func (handler *ContactHandler) getById(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil || id < 1 {
		handler.response.notFoundResponse(w, r)
		return
	}

	contact, err := handler.contactUseCase.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			handler.response.notFoundResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"contact": contact}, nil)
	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}
}

func (handler *ContactHandler) create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}

	err := helpers.ReadJSON(w, r, &input)

	if err != nil {
		handler.response.badRequestResponse(w, r, err)
		return
	}

	contact := &domain.Contact{
		FullName: input.FullName,
		Phone:    input.Phone,
	}

	v := validator.New()

	if domain.ValidateContact(v, contact); !v.Valid() {
		handler.response.failedValidationResponse(w, r, v.Errors)
	}

	err = handler.contactUseCase.Create(contact)

	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/contacts/%d", contact.ID))

	err = writeJSON(w, http.StatusCreated, envelope{"contact": contact}, headers)
	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}
}

func (handler *ContactHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)

	if err != nil || id < 1 {
		handler.response.notFoundResponse(w, r)
		return
	}

	err = handler.contactUseCase.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			handler.response.notFoundResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
	}

	err = writeJSON(w, http.StatusOK, envelope{"message": "contact successfully deleted"}, nil)
	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}
}

func (handler *ContactHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil || id < 1 {
		handler.response.notFoundResponse(w, r)
		return
	}

	contact, err := handler.contactUseCase.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			handler.response.notFoundResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FullName *string `json:"full_name"`
		Phone    *string `json:"phone"`
	}

	err = helpers.ReadJSON(w, r, &input)

	if err != nil {
		handler.response.badRequestResponse(w, r, err)
		return
	}

	if input.FullName != nil {
		contact.FullName = *input.FullName
	}

	if input.Phone != nil {
		contact.Phone = *input.Phone
	}
	v := validator.New()

	if domain.ValidateContact(v, contact); !v.Valid() {
		handler.response.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.contactUseCase.Update(contact)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			handler.response.editConflictResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
	}

	err = writeJSON(w, http.StatusOK, envelope{"contact": contact}, nil)

	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}

}
