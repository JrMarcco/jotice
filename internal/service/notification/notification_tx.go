package notification

//go:generate mockgen -source=./types.go -destination=./mock/tx_service.mock.go -package=notificationmock -type=TxService
type TxService interface{}
