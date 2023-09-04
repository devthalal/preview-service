package list_checks

import (
	common "ab-preview-service/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func ListServiceChecks() common.FunctionReturn {
	params := url.Values{}

	listAPI := os.Getenv("CONSUL_BASE_URL") + "/agent/checks"

	resp, err := http.Get(listAPI + "?" + params.Encode())
	if err != nil {
		return common.FunctionReturn{
			Message: "Request failed to get service data",
			Err:     err,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var serviceCheckMap map[string]ServiceCheckResponseData

	umErr := json.Unmarshal([]byte(body), &serviceCheckMap)
	if umErr != nil {
		fmt.Println("Unmarshal Error:", umErr)
		return common.FunctionReturn{
			Message: "Unmarshal Error getting service data",
			Err:     err,
		}
	}

	var serviceArray []ServiceCheckResponseData
	for _, service := range serviceCheckMap {
		serviceArray = append(serviceArray, service)
	}

	return common.FunctionReturn{
		Message: "Data retrieved successfully",
		Data:    serviceArray,
	}

}

func CheckLiveServices() common.FunctionReturn {
	listRes := ListServiceChecks()
	if listRes.Err != nil {
		return listRes
	}

	var liveServices []ServiceCheckResponseData

	for _, service := range listRes.Data.([]ServiceCheckResponseData) {
		if service.Status == "success" {
			liveServices = append(liveServices, service)
		}
	}

	return common.FunctionReturn{
		Message: "Data retrieved successfully",
		Data:    liveServices,
		Status:  200,
	}
}

func ListServiceStatus() common.FunctionReturn {
	checksRes := ListServiceChecks()
	if checksRes.Err != nil {
		return checksRes
	}

	var servicesData []ServiceCheckResponseData

	for _, service := range checksRes.Data.([]ServiceCheckResponseData) {
		servicesData = append(servicesData, ServiceCheckResponseData{
			Name:        service.Name,
			Status:      service.Status,
			ServiceID:   service.ServiceID,
			ServiceName: service.ServiceName,
			Output:      service.Output,
		})
	}

	return common.FunctionReturn{
		Message: "Data retrieved successfully",
		Data:    servicesData,
		Status:  200,
	}
}
