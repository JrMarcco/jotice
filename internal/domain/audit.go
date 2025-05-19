package domain

type ResourceType string

const (
	ResourceTypeTemplate ResourceType = "Template"
)

func (r ResourceType) IsTemplate() bool {
	return r == ResourceTypeTemplate
}

type Audit struct {
	ResourceId   uint64
	ResourceType ResourceType
	Content      string // json
}

type AuditContent struct {
	OwnerId       uint64   `json:"owner_id"`
	OwnerType     string   `json:"owner_type"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Channel       string   `json:"channel"`
	BizType       string   `json:"biz_type"`
	Version       string   `json:"version"`
	Signature     string   `json:"signature"`
	Content       string   `json:"content"`
	Remark        string   `json:"remark"`
	ProviderNames []string `json:"provider_names"`
}
