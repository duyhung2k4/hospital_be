package controller

import (
	"app/constant"
	"app/dto/request"
	"app/model"
	"app/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type scheduleController struct {
	query service.QueryService[model.Schedule]
}

type ScheduleController interface {
	Query(w http.ResponseWriter, r *http.Request)
}

func (c *scheduleController) Query(w http.ResponseWriter, r *http.Request) {
	var payload request.QueryReq[model.Schedule]
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	var result interface{}
	var errHandle error

	switch payload.Method {
	case constant.GET:
		result, errHandle = c.query.First(
			payload.Preload,
			payload.Omit,
			payload.Condition,
			payload.Args...,
		)
	case constant.GET_ALL:
		result, errHandle = c.query.Find(
			payload.Preload,
			payload.Omit,
			payload.Condition,
			payload.Args...,
		)
	case constant.CREATE:
		code, _ := uuid.NewV7()
		payload.Data.Code = code.String()
		payload.Data.Status = model.S_PENDING
		result, errHandle = c.query.Create(payload.Data)
	case constant.UPDATE:
		result, errHandle = c.query.Update(
			payload.Data,
			payload.Condition,
			payload.Args...,
		)
	case constant.DELETE:
		result = nil
		errHandle = c.query.Delete(
			payload.Condition,
			payload.Args...,
		)
	}

	if errHandle != nil {
		internalServerError(w, r, errHandle)
		return
	}

	res := Response{
		Data:    result,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewScheduleController() ScheduleController {
	return &scheduleController{
		query: service.NewQueryService[model.Schedule](),
	}
}
