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
