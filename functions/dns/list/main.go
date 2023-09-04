package list_dns

import (
	common "ab-preview-service/common"
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func ListDNSRecord() common.FunctionReturn {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	api, err := cloudflare.NewWithAPIToken(apiToken)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error setting cloudflare api",
			Err:     err,
			Status:  400,
		}
	}

	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	resourceContainer := cloudflare.ResourceContainer{
		Level:      "zone",
		Identifier: zoneID,
		Type:       "zone",
	}

	dNSRecords, ResultInfo, err := api.ListDNSRecords(context.Background(), &resourceContainer, cloudflare.ListDNSRecordsParams{})

	if err != nil {
		return common.FunctionReturn{
			Err:     err,
			Message: "Error getting ListDNSRecords",
		}
	}

	for _, d := range dNSRecords {
		fmt.Printf(" \n\n Name = %v ", d.Name)
		fmt.Printf(" \n ID = %v ", d.ID)
		fmt.Printf(" \n Content = %v \n\n", d.Content)
	}

	fmt.Printf(" \n ResultInfo %v \n", ResultInfo)

	return common.FunctionReturn{
		Message: "Successfully listed dns record",
		Data:    dNSRecords,
	}
}
