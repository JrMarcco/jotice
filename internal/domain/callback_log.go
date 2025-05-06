package domain

type CallbackStatus string

const (
	CallbackStatusInit    CallbackStatus = "init"
	CallbackStatusPending CallbackStatus = "pending"
	CallbackStatusSucceed CallbackStatus = "succeed"
	CallbackStatusFailed  CallbackStatus = "failed"
)

func (s CallbackStatus) String() string {
	return string(s)
}

type CallbackLog struct {
	Id           uint64
	Notification Notification
	RetryTimes   int32
	NextRetryAt  int64
	Status       CallbackStatus
}
