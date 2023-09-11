package main

import (
	"ab-preview-service/common"
	add_caddy_routes "ab-preview-service/functions/caddy/add_routes"
	"ab-preview-service/functions/caddy/create_setup"
	register_services "ab-preview-service/functions/consul/services/register"
	"ab-preview-service/functions/pull_code"
	"ab-preview-service/functions/start_package"
	"ab-preview-service/functions/unused_port"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Printf("*************************************************************")
	log.Printf("***************** Start Preview Service *********************")

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file: %v", envErr)
	}

	packageCodeUrl := os.Getenv("PACKAGE_CODE_URL")
	previewName := os.Getenv("PACKAGE_PREVIEW_NAME")
	servicesJson := os.Getenv("PACKAGE_SERVICES")

	log.Printf("**************** %v *******************", previewName)
	log.Printf("*************************************************************")

	// cRes := read_docker.ReadComposeServices()
	// cRes := delete_dns.DeleteDNSRecord([]string{})
	// cRes := list_checks.ListServiceChecks()
	// if len(cRes.Message) > 0 {
	// 	log.Printf("\n ListServiceChecks : \n %v \n", cRes.Data)
	// 	cRes = list_checks.ListServiceStatus()
	// 	log.Printf("\n ListServiceStatus : \n %v \n", cRes.Data)
	// 	cRes = list_checks.CheckLiveServices()
	// 	log.Printf("\n CheckLiveServices : \n %v \n", cRes.Data)
	// 	return
	// }

	// NOTE: our ec2 instance will be already setup with consul, caddy, docker and so
	var res common.FunctionReturn

	log.Printf("-----------Pulling code ---------------------")
	res = pull_code.PullCode(previewName, packageCodeUrl)
	log.Printf(res.Message)
	if res.Err != nil {
		log.Fatal(res.Err)
		return
	}

	var services []common.PackageServices
	log.Printf("-----------setting services ---------------------")
	if err := json.Unmarshal([]byte(servicesJson), &services); err != nil {
		log.Printf("Error Unmarshal services json ")
		log.Fatal(err)
		return
	}

	log.Printf("-----------getting unused ---------------------")
	res = unused_port.GetUnusedPort(len(services))
	log.Printf(res.Message)
	if res.Err != nil {
		log.Fatal(res.Err)
		return
	}
	fmt.Printf(" \n start_package %v \n", res.Data)

	log.Printf("-----------starting package ---------------------")
	res = start_package.StartPackage(previewName, res.Data.([]int))
	log.Printf(res.Message)
	if res.Err != nil {
		log.Fatal(res.Err)
		return
	}
	fmt.Printf(" \n start_package %v \n", res.Data)

	// log.Printf("-----------registering services ---------------------")
	domainServices := []common.PackageServices{}
	for _, service := range services {
		// register service
		res = register_services.RegisterService(service)
		log.Printf(res.Message)
		log.Printf("Data %v", res.Data)
		if res.Err != nil {
			log.Fatal(res.Err)
			return
		}

		// filter domain services
		if service.Domain != "" {
			domainServices = append(domainServices, service)
		}
	}

	// tnx code
	// res = consul_tnx.ConsulTransaction(services)
	// if res.Err != nil {
	// 	log.Printf(res.Message)
	// 	log.Fatal(res.Err)
	// 	return
	// }
	// fmt.Printf(" \n consul_tnx %v \n", res.Data)

	// check the health
	// - if it fails, deregister and mark as something wrong with start

	log.Printf("-----------setting up caddy root ---------------------")
	res = create_setup.AddCaddyConfig()
	log.Printf(res.Message)
	if res.Err != nil {
		log.Fatal(res.Err)
		return
	}
	fmt.Printf(" \n AddCaddyConfig %v \n", res.Data)

	// setup sub-routes and subject for new domain in caddy
	log.Printf("-----------setting up caddy routes ---------------------")
	res = add_caddy_routes.AddCaddyRoutes(domainServices)
	log.Printf(res.Message)
	if res.Err != nil {
		log.Fatal(res.Err)
		return
	}
	fmt.Printf(" \n AddCaddyRoutes %v \n", res.Data)

	log.Printf("#############################################################")
	log.Printf("####### End Preview Service  %v #######", previewName)
	log.Printf("#############################################################")

}

// scp -o StrictHostKeyChecking=no -i /home/ntpl/NEOITO_PROJECTS/APPBLOCKS/admin-api/_ab_previewKey.pem /home/ntpl/NEOITO_PROJECTS/APPBLOCKS/ab-preview-service/ab-preview-service ubuntu@ec2-34-229-111-40.compute-1.amazonaws.com:/home/ubuntu
