package delete_dns

import (
	common "ab-preview-service/common"
	"context"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func DeleteDNSRecord(dnsIDs []string) common.FunctionReturn {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	api, err := cloudflare.NewWithAPIToken(apiToken)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error setting cloudflare api",
			Err:     err,
		}
	}

	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	resourceContainer := cloudflare.ResourceContainer{
		Level:      "zone",
		Identifier: zoneID,
		Type:       "zone",
	}

	for _, dnsID := range dnsIDs {
		dnsErr := api.DeleteDNSRecord(context.Background(), &resourceContainer, dnsID)
		if dnsErr != nil {
			log.Printf(dnsErr.Error())
			return common.FunctionReturn{
				Message: "Error creating dns record",
				Err:     dnsErr,
			}
		}
	}

	return common.FunctionReturn{
		Message: "Successfully deleted dns record",
	}
}
