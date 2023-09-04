package register_services

type ServiceData struct {
	ID    string       `json:"ID"`
	Name  string       `json:"Name"`
	Tags  []string     `json:"Tags"`
	Port  int          `json:"Port"`
	Check ServiceCheck `json:"Check"`
	// Address string       `json:"Address"`
	// Meta    interface{}  `json:"Meta"`
}

type ServiceCheck struct {
	Http     string `json:"Http"`
	Interval string `json:"Interval"`
	Timeout  string `json:"Timeout"`
}
