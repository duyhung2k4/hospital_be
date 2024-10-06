package request

type TransitReq struct {
	ClinId        uint   `json:"clinId"`
	ScheduleId    uint   `json:"scheduleId"`
	Description   string `json:"description"`
	DepartmentIds []int  `json:"departmentIds"`
}

type SaveStepReq struct {
	ScheduleId uint   `json:"scheduleId"`
	Result     string `json:"result"`
	RoomId     uint   `json:"roomId"`
}
