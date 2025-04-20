package domain

import "github.com/JrMarcco/jotice/internal/pkg/retry"

type BizConfig struct {
	Id                uint64
	OwnerId           uint64
	OwnerType         string
	ChannelConfig     *ChannelConfig
	TransactionConfig *TransactionConfig
	RateLimit         int32
	QuotaConfig       *QuotaConfig
	CallbackConfig    *CallbackConfig
	CreateAt          int64
	UpdateAt          int64
}

type ChannelConfig struct {
	Channels    []ChannelItem `json:"channels"`
	RetryPolicy string        `json:"retry_policy"`
}

type ChannelItem struct {
	Channel  string `json:"channel"`
	Priority int32  `json:"priority"`
	Enabled  bool   `json:"enabled"`
}

type TransactionConfig struct {
	ServiceName  string        `json:"service_name"`
	InitialDelay int32         `json:"initial_delay"`
	RetryPolicy  *retry.Config `json:"retry_policy"`
}

type QuotaConfig struct {
	Daily   *DailyQuotaConfig   `json:"daily"`
	Monthly *MonthlyQuotaConfig `json:"monthly"`
}

type DailyQuotaConfig struct {
	SMS   int32 `json:"sms"`
	Email int32 `json:"email"`
}

type MonthlyQuotaConfig struct {
	SMS   int32 `json:"sms"`
	Email int32 `json:"email"`
}

type CallbackConfig struct {
	ServiceName string        `json:"service_name"`
	RetryPolicy *retry.Config `json:"retry_policy"`
}
