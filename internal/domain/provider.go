package domain

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelApp   Channel = "app"
)

func (c Channel) String() string {
	return string(c)
}

type ProviderStatus string

const (
	ProviderStatusActive   ProviderStatus = "active"
	ProviderStatusInactive ProviderStatus = "inactive"
)

func (s ProviderStatus) String() string {
	return string(s)
}

type Provider struct {
	Id      uint64
	Name    string
	Channel Channel

	Endpoint string
	RegionId string

	AppId     string
	AppKey    string
	AppSecret string

	Weight     int32
	QpsLimit   int32
	DailyLimit int32

	CallbackUrl string

	Status ProviderStatus
}
