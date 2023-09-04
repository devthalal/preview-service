package read_docker

import (
	"ab-preview-service/common"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadComposeServices(composeFilePath string) common.FunctionReturn {
	data, err := os.ReadFile(composeFilePath)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error reading file",
			Err:     err,
		}
	}

	var composeConfig DockerCompose
	err = yaml.Unmarshal(data, &composeConfig)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error parsing YAML",
			Err:     err,
		}
	}

	fmt.Println("Services in the Docker Compose file:")
	serviceNames := []string{}

	for serviceName := range composeConfig.Services {
		serviceNames = append(serviceNames, serviceName)
	}

	return common.FunctionReturn{
		Message: "Services in compose",
		Data:    serviceNames,
	}
}
