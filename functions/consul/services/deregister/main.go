package deregister_services

import (
	common "ab-preview-service/common"
	"bytes"
	"io"
	"net/http"
	"os"
)

func DeRegisterService(serviceID string) common.FunctionReturn {
	servicePutAPI := os.Getenv("CONSUL_BASE_URL") + "/agent/service/deregister/" + serviceID

	// initialize http client
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, servicePutAPI, bytes.NewBuffer([]byte{}))
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in service deregistered api setup",
			Err:     err,
			Status:  400,
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in service deregistered api ",
			Err:     err,
			Status:  400,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error reading service deregistered api response ",
			Err:     err,
			Status:  400,
		}
	}

	return common.FunctionReturn{
		Message: "Service deregistered successfully",
		Data:    body,
		Status:  200,
	}
}
