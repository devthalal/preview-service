package consul_tnx

type ServiceData struct {
	ID      string       `json:"ID"`
	Name    string       `json:"Name"`
	Tags    []string     `json:"Tags"`
	Address string       `json:"Address"`
	Port    int          `json:"Port"`
	Meta    interface{}  `json:"Meta"`
	Check   ServiceCheck `json:"Check"`
}

type ServiceCheck struct {
	ID        string `json:"ID"`
	Name      string `json:"Name"`
	ServiceID string `json:"ServiceID"`
	Http      string `json:"Http"`
	Interval  string `json:"Interval"`
	Timeout   string `json:"Timeout"`
	// DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	// TLSSkipVerify                  string `json:"TLSSkipVerify"`
}
