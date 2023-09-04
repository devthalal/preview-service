package list_services

type ServiceResponseData struct {
	ID              string      `json:"ID"`
	Service         string      `json:"Service"`
	Tags            []string    `json:"Tags"`
	Address         string      `json:"Address"`
	Port            int         `json:"Port"`
	Meta            interface{} `json:"Meta"`
	TaggedAddresses interface{} `json:"TaggedAddresses"`
}
