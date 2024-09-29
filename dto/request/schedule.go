package request

type TransitReq struct {
	Description   string `json:"description"`
	DepartmentIds []int  `json:"departmentIds"`
	ScheduleId    uint   `json:"scheduleId"`
}
