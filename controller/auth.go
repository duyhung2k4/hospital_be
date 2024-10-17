package controller

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/model"
	"app/service"
	"app/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type authController struct {
	jwtUtils     utils.JwtUtils
	authUtils    utils.AuthUtils
	queryService service.QueryService[model.Profile]
	rabbitmq     *amqp091.Connection
	redisClient  *redis.Client
	authService  service.AuthService
}

type AuthControlle interface {
	Login(w http.ResponseWriter, r *http.Request)
	CreateAdmin(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)

	Register(w http.ResponseWriter, r *http.Request)
	SendFileAuth(w http.ResponseWriter, r *http.Request)
	AuthFace(w http.ResponseWriter, r *http.Request)
	CreateSocketAuthFace(w http.ResponseWriter, r *http.Request)
	SaveProcess(w http.ResponseWriter, r *http.Request)
}

func (c *authController) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	password, _ := c.authUtils.HashPassword("tk_admin@123456")
	username := "tk_admin"

	if err := c.queryService.Delete("username = ?", username); err != nil {
		internalServerError(w, r, err)
		return
	}

	newProfile, err := c.queryService.Create(model.Profile{
		Username: username,
		Password: password,
		Role:     "admin",
		Active:   true,
	})

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	res := Response{
		Data: map[string]interface{}{
			"username": newProfile.Username,
			"password": "tk_admin@123456",
		},
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var payload request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	profileResponse, err := c.queryService.First(request.FirstPayload{Condition: "username = ?", Preload: []string{"Room"}}, payload.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	compare := c.authUtils.CheckPasswordHash(payload.Password, profileResponse.Password)
	if !compare {
		internalServerError(w, r, errors.New("password wrong"))
		return
	}

	mapData := map[string]interface{}{
		"profile_id": profileResponse.ID,
		"room_id":    profileResponse.RoomId,
		"role":       profileResponse.Role,
	}

	accessData := mapData
	accessData["uuid"] = uuid.New()
	accessData["exp"] = time.Now().Add(3 * time.Hour).Unix()
	accessToken, errAccessToken := c.jwtUtils.JwtEncode(accessData)
	if errAccessToken != nil {
		internalServerError(w, r, errAccessToken)
		return
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, errRefreshToken := c.jwtUtils.JwtEncode(refreshData)
	if errRefreshToken != nil {
		internalServerError(w, r, errRefreshToken)
		return
	}

	// errSetKeyAccessToken := a.rdb.Set(context.Background(), "access_token:"+strconv.Itoa(int(profileResponse.ID)), accessToken, 24*time.Hour).Err()
	// if errSetKeyAccessToken != nil {
	// 	internalServerError(w, r, errSetKeyAccessToken)
	// 	return
	// }
	// errSetKeyRefreshToken := a.rdb.Set(context.Background(), "refresh_token:"+strconv.Itoa(int(profileResponse.ID)), refreshToken, 3*24*time.Hour).Err()
	// if errSetKeyRefreshToken != nil {
	// 	internalServerError(w, r, errSetKeyRefreshToken)
	// 	return
	// }

	profileResponse.Password = ""
	// profileResponse.Username = ""

	res := Response{
		Data: map[string]interface{}{
			"profile":      profileResponse,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	}

	render.JSON(w, r, res)
}

func (c *authController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapDataRequest, errMapData := c.jwtUtils.JwtDecode(tokenString)

	if errMapData != nil {
		internalServerError(w, r, errMapData)
		return
	}

	profileId := uint(mapDataRequest["profile_id"].(float64))
	profileResponse, errProfile := c.queryService.First(request.FirstPayload{
		Condition: "id = ?",
		Preload:   []string{"Room"},
	}, profileId)
	if errProfile != nil {
		internalServerError(w, r, errProfile)
		return
	}

	mapData := map[string]interface{}{
		"profile_id": profileResponse.ID,
		"room_id":    profileResponse.RoomId,
		"role":       profileResponse.Role,
	}

	accessData := mapData
	accessData["uuid"] = uuid.New()
	accessData["exp"] = time.Now().Add(3 * time.Hour).Unix()
	accessToken, errAccessToken := c.jwtUtils.JwtEncode(accessData)
	if errAccessToken != nil {
		internalServerError(w, r, errAccessToken)
		return
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, errRefreshToken := c.jwtUtils.JwtEncode(refreshData)
	if errRefreshToken != nil {
		internalServerError(w, r, errRefreshToken)
		return
	}

	// errSetKeyAccessToken := a.rdb.Set(context.Background(), "access_token:"+strconv.Itoa(int(profileResponse.ID)), accessToken, 24*time.Hour).Err()
	// if errSetKeyAccessToken != nil {
	// 	internalServerError(w, r, errSetKeyAccessToken)
	// 	return
	// }
	// errSetKeyRefreshToken := a.rdb.Set(context.Background(), "refresh_token:"+strconv.Itoa(int(profileResponse.ID)), refreshToken, 3*24*time.Hour).Err()
	// if errSetKeyRefreshToken != nil {
	// 	internalServerError(w, r, errSetKeyRefreshToken)
	// 	return
	// }

	profileResponse.Password = ""
	// profileResponse.Username = ""

	res := Response{
		Data: map[string]interface{}{
			"profile":      profileResponse,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	}

	render.JSON(w, r, res)

}

func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq request.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		badRequest(w, r, err)
		return
	}

	newFolderPending := fmt.Sprintf("file/pending_file/%d", registerReq.ProfileId)
	os.RemoveAll(newFolderPending)
	if err := os.Mkdir(newFolderPending, 0777); err != nil {
		internalServerError(w, r, err)
		return
	}
	newFolderAddModel := fmt.Sprintf("file/file_add_model/%d", registerReq.ProfileId)
	os.RemoveAll(newFolderAddModel)
	if err := os.Mkdir(newFolderAddModel, 0777); err != nil {
		internalServerError(w, r, err)
		return
	}

	uuid := uuid.New().String()

	res := Response{
		Data:    uuid,
		Message: "OK",
		Error:   nil,
		Status:  200,
	}

	render.JSON(w, r, res)
}

func (c *authController) SendFileAuth(w http.ResponseWriter, r *http.Request) {
	var fileReq request.SendFileAuthFaceReq

	err := json.NewDecoder(r.Body).Decode(&fileReq)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	if fileReq.ProfileId == 0 {
		badRequest(w, r, errors.New("not found profileId"))
		return
	}

	ch, err := c.rabbitmq.Channel()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataMess := queuepayload.SendFileAuthMess{
		Data:      fileReq.Data,
		ProfileId: fileReq.ProfileId,
		Uuid:      fileReq.Uuid,
	}

	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.SEND_FILE_AUTH_QUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataMessString),
		},
	)

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

func (c *authController) AuthFace(w http.ResponseWriter, r *http.Request) {
	var authFaceReq request.AuthFaceReq
	if err := json.NewDecoder(r.Body).Decode(&authFaceReq); err != nil {
		badRequest(w, r, err)
		return
	}

	uuid := strings.Split(r.Header.Get("authorization"), " ")[1]
	if len(uuid) == 0 {
		badRequest(w, r, errors.New("not found uuid"))
		return
	}

	path, err := c.authService.CreateFileAuthFace(authFaceReq)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	ch, err := c.rabbitmq.Channel()
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	dataMess := queuepayload.FaceAuth{
		FilePath: path,
		Uuid:     uuid,
	}

	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.FACE_AUTH_QUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataMessString),
		},
	)

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

func (c *authController) CreateSocketAuthFace(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New().String()

	res := Response{
		Data:    uuid,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func (c *authController) SaveProcess(w http.ResponseWriter, r *http.Request) {
	var payload request.SaveProcessReq

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	err := c.authService.SaveFileAuth(payload.ProfileId)

	if err != nil {
		internalServerError(w, r, err)
		return
	}

	c.queryService.Update(
		model.Profile{
			Active: true,
		},
		[]string{},
		map[string][]string{},
		"id = ?",
		[]interface{}{payload.ProfileId},
	)

	res := Response{
		Data:    nil,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewAuthController() AuthControlle {
	return &authController{
		jwtUtils:     utils.NewJwtUtils(),
		authUtils:    utils.NewAuthUtils(),
		queryService: service.NewQueryService[model.Profile](),
		rabbitmq:     config.GetRabbitmq(),
		redisClient:  config.GetRedisClient(),
		authService:  service.NewAuthService(),
	}
}
