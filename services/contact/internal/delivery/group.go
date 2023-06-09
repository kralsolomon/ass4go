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

type GroupHandler struct {
	groupUseCase domain.GroupUseCase
	response     responseHandler
}

func NewGroupHandler(router *httprouter.Router, logger *jsonlog.Logger, groupUseCase domain.GroupUseCase) {
	handler := &GroupHandler{
		groupUseCase: groupUseCase,
		response:     responseHandler{logger: logger},
	}
	router.HandlerFunc(http.MethodGet, "/group", handler.getById)
	router.HandlerFunc(http.MethodPost, "/group", handler.create)
	router.HandlerFunc(http.MethodPut, "/group", handler.update)
	router.HandlerFunc(http.MethodGet, "/group/healthcheck", handler.healthcheck)
}

func (handler *GroupHandler) healthcheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, envelope{"status": "ok"}, nil)
}
func (handler *GroupHandler) getById(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil || id < 1 {
		handler.response.notFoundResponse(w, r)
		return
	}

	group, err := handler.groupUseCase.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			handler.response.notFoundResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"group": group}, nil)
	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}
}

func (handler *GroupHandler) create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GroupName string `json:"group_name"`
	}

	err := helpers.ReadJSON(w, r, &input)

	if err != nil {
		handler.response.badRequestResponse(w, r, err)
		return
	}

	group := &domain.Group{
		GroupName: input.GroupName,
	}

	v := validator.New()

	if domain.ValidateGroup(v, group); !v.Valid() {
		handler.response.failedValidationResponse(w, r, v.Errors)
	}

	err = handler.groupUseCase.Create(group)

	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/group/%d", group.ID))

	err = writeJSON(w, http.StatusCreated, envelope{"group": group}, headers)
	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}
}

func (handler *GroupHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil || id < 1 {
		handler.response.notFoundResponse(w, r)
		return
	}

	group, err := handler.groupUseCase.GetByID(id)
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
		GroupName *string `json:"group_name"`
	}

	err = helpers.ReadJSON(w, r, &input)

	if err != nil {
		handler.response.badRequestResponse(w, r, err)
		return
	}

	if input.GroupName != nil {
		group.GroupName = *input.GroupName
	}

	v := validator.New()

	if domain.ValidateGroup(v, group); !v.Valid() {
		handler.response.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.groupUseCase.Update(group)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			handler.response.editConflictResponse(w, r)
		default:
			handler.response.serverErrorResponse(w, r, err)
		}
	}

	err = writeJSON(w, http.StatusOK, envelope{"group": group}, nil)

	if err != nil {
		handler.response.serverErrorResponse(w, r, err)
	}

}
