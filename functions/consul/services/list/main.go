package list_services

import (
	common "ab-preview-service/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func ListServices() common.FunctionReturn {
	params := url.Values{}

	listAPI := os.Getenv("CONSUL_BASE_URL") + "/agent/services"

	resp, err := http.Get(listAPI + "?" + params.Encode())
	if err != nil {
		return common.FunctionReturn{
			Message: "Request failed to get service data",
			Err:     err,
			Status:  400,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var serviceMap map[string]ServiceResponseData

	umErr := json.Unmarshal([]byte(body), &serviceMap)
	if umErr != nil {
		fmt.Println("Unmarshal Error:", umErr)
		return common.FunctionReturn{
			Message: "Unmarshal Error getting service data",
			Err:     err,
			Status:  400,
		}
	}

	return common.FunctionReturn{
		Message: "Data retrieved successfully",
		Data:    serviceMap,
		Status:  200,
	}

}

func CheckLiveServices() {}
