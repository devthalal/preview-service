package create_setup

import (
	common "ab-preview-service/common"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func AddCaddyConfig() common.FunctionReturn {
	caddyConfigAPI := os.Getenv("CADDY_BASE_URL") + "/load"
	config := map[string]interface{}{
		"apps": map[string]interface{}{
			"http": map[string]interface{}{
				"servers": map[string]interface{}{
					"preview_server": map[string]interface{}{
						"listen": []string{":443", ":80"},
						"routes": []map[string]interface{}{},
					},
				},
			},
			"tls": map[string]interface{}{
				"automation": map[string]interface{}{
					"policies": []map[string]interface{}{
						{
							"issuers": []map[string]interface{}{
								{"email": "info@appblocks.com", "module": "acme"},
								{"email": "info@appblocks.com", "module": "zerossl"},
							},
							"subjects": []string{},
						},
					},
				},
			},
		},
	}

	// initialize http client
	client := &http.Client{}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return common.FunctionReturn{
			Message: "Error in marshaling caddy JSON data",
			Err:     err,
		}
	}

	req, err := http.NewRequest(http.MethodPost, caddyConfigAPI, bytes.NewBuffer(jsonData))
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

	return common.FunctionReturn{
		Message: "Caddy added route successfully",
		Data:    string(body),
		Status:  200,
	}
}
