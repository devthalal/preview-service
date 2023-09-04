package register_services

import (
	common "ab-preview-service/common"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func RegisterService(service common.PackageServices) common.FunctionReturn {
	servicePutAPI := os.Getenv("CONSUL_BASE_URL") + "/agent/service/register?replace-existing-checks=true"
	serviceData := ServiceData{
		ID:   service.Name,
		Name: service.Name,
		Tags: service.Tags,
		Port: service.Port,
		Check: ServiceCheck{
			Http:     "http://localhost:" + strconv.Itoa(service.Port),
			Interval: "30s",
			Timeout:  "5s",
		},
	}
	serviceDataJson, err := json.Marshal(serviceData)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error formatting service data",
			Err:     err,
			Status:  400,
		}
	}
	log.Printf("serviceDataJson %v", string(serviceDataJson))

	// initialize http client
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, servicePutAPI, bytes.NewBuffer(serviceDataJson))
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in service register api ",
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
