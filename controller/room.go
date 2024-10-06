package controller

import (
	"app/dto/request"
	"app/service"
	"app/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type roomController struct {
	roomService service.RoomService
	stepService service.StepService
	jwtUtils    utils.JwtUtils
}

type RoomController interface {
	SaveStep(w http.ResponseWriter, r *http.Request)
	CallStep(w http.ResponseWriter, r *http.Request)
	PullStep(w http.ResponseWriter, r *http.Request)
	AddAccount(w http.ResponseWriter, r *http.Request)
}

func (c *roomController) SaveStep(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapData, _ := c.jwtUtils.JwtDecode(tokenString)
	if mapData["room_id"] == nil {
		internalServerError(w, r, errors.New("room_id not found"))
		return
	}
	roomId := uint(mapData["room_id"].(float64))

	var payload request.SaveStepReq
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}
	payload.RoomId = roomId

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
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapData, _ := c.jwtUtils.JwtDecode(tokenString)
	if mapData["room_id"] == nil {
		internalServerError(w, r, errors.New("room_id not found"))
		return
	}
	roomId := uint(mapData["room_id"].(float64))

	step, err := c.roomService.CallStep(roomId)

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
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapData, _ := c.jwtUtils.JwtDecode(tokenString)
	if mapData["room_id"] == nil {
		internalServerError(w, r, errors.New("room_id not found"))
		return
	}
	roomId := uint(mapData["room_id"].(float64))

	step, err := c.roomService.PullStep(roomId)

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
		jwtUtils:    utils.NewJwtUtils(),
	}
}
