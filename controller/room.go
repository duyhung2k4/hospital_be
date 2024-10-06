package controller

import (
	"app/dto/request"
	"app/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

type roomController struct {
	roomService service.RoomService
	stepService service.StepService
}

type RoomController interface {
	SaveStep(w http.ResponseWriter, r *http.Request)
	CallStep(w http.ResponseWriter, r *http.Request)
	PullStep(w http.ResponseWriter, r *http.Request)
	AddAccount(w http.ResponseWriter, r *http.Request)
}

func (c *roomController) SaveStep(w http.ResponseWriter, r *http.Request) {
	var payload request.SaveStepReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	payload.RoomId = 3

	err := c.stepService.SaveStep(payload)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = c.stepService.NextStep(payload.ScheduleId)
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

func (c *roomController) CallStep(w http.ResponseWriter, r *http.Request) {
	step, err := c.roomService.CallStep(3)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    step,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *roomController) PullStep(w http.ResponseWriter, r *http.Request) {
	step, err := c.roomService.PullStep(3)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    step,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *roomController) AddAccount(w http.ResponseWriter, r *http.Request) {
	var payload request.AddAccountRoomReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	room, err := c.roomService.AddAccount(payload)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data:    room,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewRoomController() RoomController {
	return &roomController{
		roomService: service.NewRoomService(),
		stepService: service.NewStepService(),
	}
}
