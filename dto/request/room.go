package request

type AddAccountRoomReq struct {
	RoomId      uint   `json:"roomId"`
	Password    string `json:"password"`
	EmailAccept string `json:"emailAccept"`
}
