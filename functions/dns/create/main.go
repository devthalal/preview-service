package create_dns

import (
	common "ab-preview-service/common"
	delete_dns "ab-preview-service/functions/dns/delete"
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func CreateDNSRecord(domainServices []common.PackageServices) common.FunctionReturn {

	dnsData := make([]CreateDNSInput, len(domainServices))
	for i, service := range domainServices {
		dnsData[i] = CreateDNSInput{
			ServerContent: os.Getenv("PREVIEW_INSTANCE_PUBLIC_DNS"),
			DomainURL:     service.Domain,
			Name:          service.Name,
		}
	}

	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

	api, err := cloudflare.NewWithAPIToken(apiToken)

	if err != nil {
		return common.FunctionReturn{
			Message: "Error setting cloudflare api",
			Err:     err,
		}
	}

	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	proxied := true

	dnsErrors := make([]common.FunctionReturn, 0)
	dnsRecords := make([]cloudflare.DNSRecord, 0)

	resourceContainer := cloudflare.ResourceContainer{
		Level:      "zone",
		Type:       "zone",
		Identifier: zoneID,
	}

	for _, data := range dnsData {

		dnsRecordParams := cloudflare.CreateDNSRecordParams{
			Type:    "CNAME",
			ZoneID:  zoneID,
			Name:    data.Name,
			Content: data.ServerContent,
			Proxied: &proxied,
		}

		dnsRes, dnsErr := api.CreateDNSRecord(context.Background(), &resourceContainer, dnsRecordParams)
		if dnsErr != nil {
			dnsErrors = append(dnsErrors, common.FunctionReturn{
				Message: "Error creating dns record",
				Err:     dnsErr,
			})
			continue
		}

		fmt.Printf("\n=== data.DomainURL %s \nName %v \nID %v \nContent %v  === \n ", data.DomainURL, dnsRes.Name, dnsRes.ID, dnsRes.Content)

		dnsRecords = append(dnsRecords, dnsRes)

	}

	if len(dnsErrors) > 0 {
		dnsIDs := []string{}
		for _, item := range dnsRecords {
			if item.ID != "" {
				dnsIDs = append(dnsIDs, item.ID)
			}
		}

		fmt.Printf("\n!!!!!!!! DELETING DNS RECORDS %v !!!!!!!!", dnsIDs)

		res := delete_dns.DeleteDNSRecord(dnsIDs)
		if res.Err != nil {
			return res
		}

		return common.FunctionReturn{
			Message: "Error creating domains",
			Data:    dnsErrors,
		}
	}

	return common.FunctionReturn{
		Message: "Domains created successfully",
		Data:    dnsRecords,
	}
}
