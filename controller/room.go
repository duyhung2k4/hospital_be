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
}

type RoomController interface {
	CallStep(w http.ResponseWriter, r *http.Request)
	PullStep(w http.ResponseWriter, r *http.Request)
	AddAccount(w http.ResponseWriter, r *http.Request)
}

func (c *roomController) CallStep(w http.ResponseWriter, r *http.Request) {
	step, err := c.roomService.CallStep(6)
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
	step, err := c.roomService.PullStep(6)
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
	}
}
