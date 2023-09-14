package add_caddy_routes

import (
	common "ab-preview-service/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func AddCaddyRoutes(services []common.PackageServices) common.FunctionReturn {
	caddyPutAPI := os.Getenv("CADDY_BASE_URL") + "/config/apps/http/servers/preview_server/routes/..."
	putData := getCaddyRouteData(services)

	// initialize http client
	client := &http.Client{}

	jsonData, err := json.Marshal(putData)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in marshaling caddy JSON data",
			Err:     err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, caddyPutAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in caddy add route api ",
			Err:     err,
			Status:  400,
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in caddy add route api ",
			Err:     err,
			Status:  400,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error reading caddy add route api response ",
			Err:     err,
			Status:  400,
		}
	}

	// Add subjects
	subjRes := AddCaddyTLSSubjects(services)

	if subjRes.Err != nil {
		return subjRes
	}

	bodyString := fmt.Sprintf("caddy  %s \n tls_subject %s", string(body), subjRes.Data)

	return common.FunctionReturn{
		Message: "Caddy added route successfully",
		Data:    bodyString,
		Status:  200,
	}
}

func getCaddyRouteData(services []common.PackageServices) []map[string]interface{} {
	routesDatas := make([]map[string]interface{}, 0)

	for _, route := range services {

		buildPath := ""

		if strings.Contains(route.Domain, "elements") {
			buildPath = "/usr/share/elements_build"
		} else {
			buildPath = "/usr/share/container_build"
		}

		routeHandle := []map[string]interface{}{
			{
				"handler": "vars",
				"root":    buildPath,
			},
			{
				"handler":     "file_server",
				"index_names": []string{"index.html"},
			},
		}

		if strings.Contains(route.Domain, "function") {
			routeHandle = []map[string]interface{}{
				{
					"handler": "headers",
					"response": map[string]interface{}{
						"set": map[string]interface{}{
							"Cache-Control": []string{"no-cache"},
						},
					},
				},
				{
					"handler": "reverse_proxy",
					"upstreams": []map[string]interface{}{
						{
							"dial": "0.0.0.0:" + strconv.Itoa(route.Port),
						},
					},
				},
			}
		}

		routeData := map[string]interface{}{
			"match": []map[string]interface{}{
				{
					"host": []string{route.Domain},
				},
			},
			"handle": []map[string]interface{}{
				{
					"handler": "subroute",
					"routes": []map[string]interface{}{
						{
							"handle": routeHandle,
						},
					},
				},
			},
			"terminal": true,
		}

		routesDatas = append(routesDatas, routeData)
	}

	return routesDatas
}

func AddCaddyTLSSubjects(services []common.PackageServices) common.FunctionReturn {
	caddyPutSubjectsAPI := os.Getenv("CADDY_BASE_URL") + "/config/apps/tls/automation/policies/0/subjects/..."

	domains := make([]string, 0)
	for _, item := range services {
		domains = append(domains, item.Domain)
	}

	putData := domains // Slice of strings
	jsonData, err := json.Marshal(putData)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in marshaling tls subject JSON data",
			Err:     err,
		}
	}

	// initialize http client
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, caddyPutSubjectsAPI, bytes.NewReader(jsonData))
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in caddy add tls subject api ",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in caddy add tls subject api ",
			Err:     err,
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error reading caddy add tls subject api response ",
			Err:     err,
		}
	}

	return common.FunctionReturn{
		Message: "Caddy added subject successfully",
		Data:    string(body),
	}
}
