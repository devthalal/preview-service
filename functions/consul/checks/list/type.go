package list_checks

type ServiceCheckResponseData struct {
	Node        string   `json:"Node"`
	Status      string   `json:"Status"`
	CheckID     string   `json:"CheckID"`
	ServiceID   string   `json:"ServiceID"`
	Name        string   `json:"Name,omitempty"`
	Notes       string   `json:"Notes,omitempty"`
	Output      string   `json:"Output,omitempty"`
	ServiceName string   `json:"ServiceName,omitempty"`
	ServiceTags []string `json:"ServiceTags,omitempty"`
	Definition  struct{} `json:"Definition,omitempty"`
	Type        string   `json:"Type,omitempty"`
	Interval    string   `json:"Interval,omitempty"`
	Timeout     string   `json:"Timeout,omitempty"`
	ExposedPort int      `json:"ExposedPort,omitempty"`
	CreateIndex int      `json:"CreateIndex,omitempty"`
	ModifyIndex int      `json:"ModifyIndex,omitempty"`
}
