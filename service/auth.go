package service

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"app/model"
	"app/utils"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authService struct {
	psql     *gorm.DB
	redis    *redis.Client
	jwtUtils utils.JwtUtils
	rabbitmq *amqp091.Connection
}

type AuthService interface {
	CheckFace(payload queuepayload.SendFileAuthMess) (string, error)
	CreateFileAuthFace(data request.AuthFaceReq) (string, error)
	AuthFace(payload queuepayload.FaceAuth) (int, error)
	ActiveProfile(auth string) (*model.Profile, error)
	SaveFileAuth(profileId uint) error
	GetProfile(profileId uint) (*model.Profile, error)
	CreateToken(profileId uint) (string, string, error)
	ShowCheck(payload queuepayload.ShowCheck) error
}

func (s *authService) CheckFace(payload queuepayload.SendFileAuthMess) (string, error) {
	base64Data := payload.Data
	imgData, err := base64.StdEncoding.DecodeString(base64Data[strings.IndexByte(base64Data, ',')+1:])
	if err != nil {
		log.Println(err)
		return "", err
	}

	fileName := uuid.New().String()

	// Check num image for train
	pathCheckNumFolder := fmt.Sprintf("file/file_add_model/%d", payload.ProfileId)
	countFileFolder, err := utils.CheckNumFolder(pathCheckNumFolder)
	if err != nil {
		return "", err
	}
	// Config num input data
	if countFileFolder >= 5 {
		return "done", nil
	}

	// Tạo file tạm thời từ dữ liệu ảnh
	pathPending := fmt.Sprintf("file/pending_file/%d/%s.png", payload.ProfileId, fileName)
	filePending, err := os.Create(pathPending)
	if err != nil {
		return "", err
	}
	defer filePending.Close()

	_, err = filePending.Write(imgData)
	if err != nil {
		return "", err
	}

	payloadDetectFace, err := json.Marshal(map[string]interface{}{
		"input_image_path": pathPending,
	})
	if err != nil {
		return "", err
	}

	// Gọi API để kiểm tra khuôn mặt
	resp, err := http.Post("http://localhost:5000/detect_single_face", "application/json", bytes.NewBuffer(payloadDetectFace))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}
	// Đọc phản hồi từ API
	var resultCheckFace struct {
		Result bool `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultCheckFace); err != nil {
		return "", err
	}
	if !resultCheckFace.Result {
		if err := os.Remove(pathPending); err != nil {
			return "", err
		}
		return "image not a face!", nil
	}

	// Gọi API tính góc xoay
	resp, err = http.Post("http://localhost:5000/calculate_head_pose", "application/json", bytes.NewBuffer(payloadDetectFace))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}
	// Đọc phản hồi từ API
	if err := json.NewDecoder(resp.Body).Decode(&resultCheckFace); err != nil {
		return "", err
	}
	if !resultCheckFace.Result {
		if err := os.Remove(pathPending); err != nil {
			return "", err
		}
		return "image not a face!", nil
	}

	// Thêm dữ liệu vào mô hình
	pathAddModel := fmt.Sprintf("file/file_add_model/%d/%s.png", payload.ProfileId, fileName)
	fileAddModel, err := os.Create(pathAddModel)
	if err != nil {
		return "", err
	}
	defer fileAddModel.Close()

	_, err = fileAddModel.Write(imgData)
	if err != nil {
		return "", err
	}

	return "not enough data", nil
}

func (s *authService) CreateFileAuthFace(data request.AuthFaceReq) (string, error) {
	base64Data := data.Data
	imgData, err := base64.StdEncoding.DecodeString(base64Data[strings.IndexByte(base64Data, ',')+1:])
	fileName := uuid.New().String()

	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("file/auth_face/%s.png", fileName)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	_, err = file.Write(imgData)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *authService) AuthFace(payload queuepayload.FaceAuth) (int, error) {
	var faces []model.Face

	// Lấy danh sách khuôn mặt từ cơ sở dữ liệu
	if err := s.psql.
		Model(&model.Face{}).
		Joins("JOIN profiles AS p ON p.id = faces.profile_id").
		Where("p.active = ?", true).
		Find(&faces).
		Error; err != nil {
		return 0, err
	}

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]interface{}{
		"faces":            faces,
		"input_image_path": payload.FilePath,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/recognize_faces", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		os.Remove(payload.FilePath)
		return 0, err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		os.Remove(payload.FilePath)
		return 0, fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var response struct {
		Result   string  `json:"result"`
		Accuracy float64 `json:"accuracy"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	// Xử lý kết quả từ API
	// if response.Result == "-1" {
	// 	return -1, nil // Không tìm thấy khuôn mặt phù hợp
	// }

	profileId, err := strconv.Atoi(response.Result)
	if err != nil {
		return 0, err
	}

	if profileId <= 0 {
		os.Remove(payload.FilePath)
		return profileId, nil
	}

	// show check
	ch, err := s.rabbitmq.Channel()
	if err != nil {
		return 0, err
	}
	dataMess := queuepayload.ShowCheck{
		FilePath:  payload.FilePath,
		ProfileId: response.Result,
		Accuracy:  response.Accuracy,
	}

	dataMessString, err := json.Marshal(dataMess)
	if err != nil {
		return 0, err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		string(constant.SHOW_CHECK_QUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataMessString),
		},
	)
	if err != nil {
		return 0, err
	}

	return profileId, nil
}

func (s *authService) ActiveProfile(auth string) (*model.Profile, error) {
	var profile model.Profile
	profileJson, err := s.redis.Get(context.Background(), auth).Result()

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(profileJson), &profile); err != nil {
		return nil, err
	}

	if err := s.psql.
		Model(&model.Profile{}).
		Where("id = ?", profile.ID).
		Updates(&model.Profile{Active: true}).
		Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *authService) SaveFileAuth(profileId uint) error {
	// Tạo đường dẫn đến thư mục chứa ảnh
	pathFileAddModel := fmt.Sprintf("file/file_add_model/%d", profileId)

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]interface{}{
		"directory_path": pathFileAddModel,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/face_encoding", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call API, status code: %d", resp.StatusCode)
	}

	// Đọc phản hồi từ API
	var response struct {
		Result        string      `json:"result"`
		FaceEncodings [][]float64 `json:"face_encodings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Result != "success" {
		return fmt.Errorf("API error: %s", response.Result)
	}

	// Thêm khuôn mặt vào danh sách
	var faces []model.Face
	for _, img := range response.FaceEncodings {
		faces = append(faces, model.Face{
			ProfileId:    profileId,
			FaceEncoding: img,
		})
	}

	if err := s.psql.Model(&model.Face{}).Create(&faces).Error; err != nil {
		return err
	}

	// Xóa thư mục tạm
	pendingPath := fmt.Sprintf("file/pending_file/%d", profileId)
	if err := os.RemoveAll(pendingPath); err != nil {
		return err
	}
	addModelPath := fmt.Sprintf("file/file_add_model/%d", profileId)
	if err := os.RemoveAll(addModelPath); err != nil {
		return err
	}

	return nil
}

func (s *authService) GetProfile(profileId uint) (*model.Profile, error) {
	var profile *model.Profile
	if err := s.psql.
		Model(&model.Profile{}).
		Where("id = ?", profileId).
		First(&profile).
		Error; err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *authService) CreateToken(profileId uint) (string, string, error) {
	var profile *model.Profile

	if err := s.psql.
		Model(&model.Profile{}).
		Where("id = ?", profileId).
		First(&profile).
		Error; err != nil {
		return "", "", err
	}

	mapData := map[string]interface{}{
		"profile_id": profile.ID,
		"email":      profile.Email,
	}

	accessData := mapData
	accessData["uuid"] = uuid.New()
	accessData["exp"] = time.Now().Add(3 * time.Hour).Unix()
	accessToken, err := s.jwtUtils.JwtEncode(accessData)
	if err != nil {
		return "", "", err
	}

	refreshData := mapData
	refreshData["uuid"] = uuid.New()
	refreshData["exp"] = time.Now().Add(3 * 3 * time.Hour).Unix()
	refreshToken, err := s.jwtUtils.JwtEncode(refreshData)
	if err != nil {
		return "", "", err
	}

	err = s.redis.Set(context.Background(), "access_token:"+strconv.Itoa(int(profile.ID)), accessToken, 24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}
	err = s.redis.Set(context.Background(), "refresh_token:"+strconv.Itoa(int(profile.ID)), refreshToken, 3*24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ShowCheck(payload queuepayload.ShowCheck) error {
	// Tạo đường dẫn đến thư mục chứa ảnh
	filename := strings.Split(payload.FilePath, "/")[2]
	pathFileAuthFace := fmt.Sprintf("file/auth_face/%s", filename)
	saveFileAuthFace := fmt.Sprintf("file/save_auth/%s", filename)

	// Tạo dữ liệu JSON để gửi đến API
	data := map[string]interface{}{
		"input_image_path": pathFileAuthFace,
		"save_path":        saveFileAuthFace,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Gửi yêu cầu POST đến API Flask
	resp, err := http.Post("http://localhost:5000/show_check", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer os.Remove(pathFileAuthFace)

	profileId, err := strconv.Atoi(payload.ProfileId)
	if err != nil {
		return err
	}

	log.Println(payload.ProfileId)
	log.Println(profileId)

	var logCheck *model.LogCheck = &model.LogCheck{
		ProfileId: uint(profileId),
		Accuracy:  payload.Accuracy,
		Url:       saveFileAuthFace,
	}

	if err := s.psql.Model(&model.LogCheck{}).Create(&logCheck).Error; err != nil {
		return err
	}

	return nil
}

func NewAuthService() AuthService {
	return &authService{
		redis:    config.GetRedisClient(),
		psql:     config.GetPsql(),
		rabbitmq: config.GetRabbitmq(),
		jwtUtils: utils.NewJwtUtils(),
	}
}
