package create_dns

type CreateDNSInput struct {
	AppID         string `json:"app_id"`
	EnvironmentID string `json:"environment_id"`
	Name          string `json:"name"`
	DomainURL     string `json:"domain_url"`
	ServerContent string `json:"server_content"`
	StaticUrl     bool   `json:"static_url"`
}
