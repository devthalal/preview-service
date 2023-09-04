package common

type FunctionReturn struct {
	Data    interface{}
	Err     interface{}
	Message string
	Status  int
}

type PackageServices struct {
	Name   string   `json:"Name"`
	Domain string   `json:"Domain"`
	Port   int      `json:"Port"`
	Tags   []string `json:"Tags"`
}
