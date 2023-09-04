package consul_tnx

import (
	common "ab-preview-service/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func ConsulTransaction(services []common.PackageServices) common.FunctionReturn {
	servicePutAPI := os.Getenv("CONSUL_BASE_URL") + "/tnx"

	serviceDataMap := make([]ServiceData, 0)
	for i, service := range services {
		serviceDataMap[i] = ServiceData{
			ID:   service.Name,
			Name: service.Name,
			Tags: service.Tags,
			Port: service.Port,
			Check: ServiceCheck{
				ID:        "check-" + service.Name,
				Name:      "check-" + service.Name,
				ServiceID: service.Name,
				Http:      "http://localhost:" + strconv.Itoa(service.Port),
				Interval:  "30s",
				Timeout:   "5s",
			},
		}
	}

	serviceDataJson, err := json.Marshal(serviceDataMap)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error formatting service data",
			Err:     err,
		}
	}

	fmt.Printf("servicePutAPI %s \n", servicePutAPI)
	fmt.Printf("serviceDataJson %v \n", string(serviceDataJson))

	// initialize http client
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, servicePutAPI, bytes.NewBuffer(serviceDataJson))
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in setting up service register api ",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in service register api ",
			Err:     err,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error reading service register api response ",
			Err:     err,
		}
	}

	return common.FunctionReturn{
		Message: "Service registered successfully",
		Data:    string(body),
	}
}
