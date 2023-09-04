package read_docker

type DockerCompose struct {
	Services map[string]interface{} `yaml:"services"`
}
