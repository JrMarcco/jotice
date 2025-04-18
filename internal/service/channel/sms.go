package channel

type smsChannel struct {
	baseChannel
}

func NewSMSChannel() Channel {
	return &smsChannel{
		baseChannel: baseChannel{},
	}
}
