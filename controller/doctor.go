package controller

type docterController struct{}

type DocterController interface {
}

func NewDocterController() DocterController {
	return &docterController{}
}
