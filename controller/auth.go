package controller

import (
	"app/dto/request"
	"app/model"
	"app/service"
	"app/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type authController struct {
	jwtUtils     utils.JwtUtils
	authUtils    utils.AuthUtils
	queryService service.QueryService[model.Profile]
}

type AuthControlle interface {
	Login(w http.ResponseWriter, r *http.Request)
	CreateAdmin(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

func (a *authController) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	password, _ := a.authUtils.HashPassword("tk_admin@123456")
	username := "tk_admin"

	if err := a.queryService.Delete("username = ?", username); err != nil {
		internalServerError(w, r, err)
		return
	}

	newProfile, err := a.queryService.Create(model.Profile{
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

func (a *authController) Login(w http.ResponseWriter, r *http.Request) {
	var payload request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		badRequest(w, r, err)
		return
	}

	profileResponse, err := a.queryService.First(request.FirstPayload{Condition: "username = ?", Preload: []string{"Room"}}, payload.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	compare := a.authUtils.CheckPasswordHash(payload.Password, profileResponse.Password)
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
	accessToken, errAccessToken := a.jwtUtils.JwtEncode(accessData)
	if errAccessToken != nil {
		internalServerError(w, r, errAccessToken)
		return
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, errRefreshToken := a.jwtUtils.JwtEncode(refreshData)
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

func (a *authController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	mapDataRequest, errMapData := a.jwtUtils.JwtDecode(tokenString)

	if errMapData != nil {
		internalServerError(w, r, errMapData)
		return
	}

	profileId := uint(mapDataRequest["profile_id"].(float64))
	profileResponse, errProfile := a.queryService.First(request.FirstPayload{
		Condition: "id = ?",
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
	accessToken, errAccessToken := a.jwtUtils.JwtEncode(accessData)
	if errAccessToken != nil {
		internalServerError(w, r, errAccessToken)
		return
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, errRefreshToken := a.jwtUtils.JwtEncode(refreshData)
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

func NewAuthController() AuthControlle {
	return &authController{
		jwtUtils:     utils.NewJwtUtils(),
		authUtils:    utils.NewAuthUtils(),
		queryService: service.NewQueryService[model.Profile](),
	}
}
