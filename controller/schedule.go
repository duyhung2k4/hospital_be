package controller

import (
	"app/constant"
	"app/dto/request"
	"app/model"
	"app/service"
	"app/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type scheduleController struct {
	query           service.QueryService[model.Schedule]
	scheduleService service.ScheduleService
	jwtUtils        utils.JwtUtils
}

type ScheduleController interface {
	Query(w http.ResponseWriter, r *http.Request)
	CallMedicalFile(w http.ResponseWriter, r *http.Request)
	Transit(w http.ResponseWriter, r *http.Request)
	PullMedicalFile(w http.ResponseWriter, r *http.Request)
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
			request.FirstPayload{
				Preload:   payload.Preload,
				Omit:      payload.Omit,
				Condition: payload.Condition,
			},
			payload.Args...,
		)
	case constant.GET_ALL:
		result, errHandle = c.query.Find(
			request.FindPayload{
				Preload:   payload.Preload,
				Omit:      payload.Omit,
				Condition: payload.Condition,
				Order:     payload.Order,
			},
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
			payload.Preload,
			payload.Omit,
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

func (c *scheduleController) CallMedicalFile(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapData, _ := c.jwtUtils.JwtDecode(tokenString)
	if mapData["room_id"] == nil {
		internalServerError(w, r, errors.New("room_id not found"))
		return
	}
	roomId := uint(mapData["room_id"].(float64))

	var payload request.QueryReq[model.Schedule] = request.QueryReq[model.Schedule]{
		Preload:   []string{},
		Omit:      map[string][]string{},
		Condition: "status = ? AND room_id = ?",
		Args:      []interface{}{model.S_EXAMINING, roomId},
	}

	schedule, err := c.query.First(request.FirstPayload{
		Preload:   payload.Preload,
		Omit:      payload.Omit,
		Condition: payload.Condition,
	}, payload.Args...)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("cc")
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    schedule,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (s *scheduleController) Transit(w http.ResponseWriter, r *http.Request) {
	var payload request.TransitReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	err := s.scheduleService.Transit(payload)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    nil,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *scheduleController) PullMedicalFile(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapData, _ := c.jwtUtils.JwtDecode(tokenString)
	if mapData["room_id"] == nil {
		internalServerError(w, r, errors.New("room_id not found"))
		return
	}
	roomId := uint(mapData["room_id"].(float64))

	schedule, err := c.scheduleService.PullSchedule(roomId)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    schedule,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewScheduleController() ScheduleController {
	return &scheduleController{
		query:           service.NewQueryService[model.Schedule](),
		scheduleService: service.NewScheduleService(),
		jwtUtils:        utils.NewJwtUtils(),
	}
}
