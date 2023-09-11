package common

type FunctionReturn struct {
	Data    interface{}
	Err     interface{}
	Message string
	Status  int
}

type PackageServices struct {
	Name          string   `json:"Name,omitempty"`
	Domain        string   `json:"Domain,omitempty"`
	Port          int      `json:"Port,omitempty"`
	ServerContent string   `json:"ServerContent,omitempty"`
	Tags          []string `json:"Tags,omitempty"`
}
