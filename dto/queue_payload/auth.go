package queuepayload

type SendFileAuthMess struct {
	ProfileId uint
	Uuid      string
	Data      string
}

type FaceAuth struct {
	Uuid     string
	FilePath string
}

type ShowCheck struct {
	ProfileId string  `json:"profileId"`
	Accuracy  float64 `json:"accuracy"`
	FilePath  string  `json:"filePath"`
}
