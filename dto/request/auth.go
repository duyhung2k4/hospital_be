package request

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	ProfileId uint `json:"profileId"`
}

type SendFileAuthFaceReq struct {
	Data      string `json:"data"`
	ProfileId uint   `json:"profileId"`
	Uuid      string `json:"uuid"`
}

type AuthFaceReq struct {
	Data string `json:"data"`
}

type AcceptCodeReq struct {
	Code string `json:"code"`
}

type SaveProcessReq struct {
	ProfileId uint `json:"profileId"`
}
