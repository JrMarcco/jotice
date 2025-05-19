package template

//go:generate mockgen -source=./types.go -destination=./mock/tpl_service.mock.go -package=templatemock -type=TplService
type TplService interface{}
