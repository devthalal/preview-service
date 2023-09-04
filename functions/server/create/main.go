package create_server

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/aidarkhanov/nanoid"
// 	"github.com/appblocks-hub/appblocks-datamodels-backend/models"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// 	"github.com/cloudflare/cloudflare-go"
// 	"gorm.io/gorm"
// )

// func CreateServer() {

// 	appId := ""
// 	var Db *gorm.DB

// 	awsSession, err := session.NewSession(&aws.Config{
// 		Region: aws.String(os.Getenv("S3_BLOCK_EC2_BUCKET_REGION"))},
// 	)

// 	if err != nil {
// 		fmt.Println("Failed to create a new session: ", err)
// 	}

// 	// Create EC2 service client
// 	svc := ec2.New(awsSession)

// 	// Create key pair .

// 	KeyName := aws.String(appId + "__" + strconv.FormatInt(time.Now().Unix(), 10))
// 	keyPairResult, err := svc.CreateKeyPair(&ec2.CreateKeyPairInput{
// 		KeyName: KeyName,
// 	})

// 	if err != nil {
// 		fmt.Println("Could not create keyPairResult", err)
// 	}

// 	// Specify the details of the instance that you want to create.

// 	runInstanceResult, err := svc.RunInstances(&ec2.RunInstancesInput{
// 		ImageId:      aws.String(os.Getenv("EC2_INSTANCE_IMAGE_ID")),
// 		InstanceType: aws.String(os.Getenv("EC2_INSTANCE_TYPE")),
// 		MinCount:     aws.Int64(1),
// 		MaxCount:     aws.Int64(1),
// 		KeyName:      KeyName,
// 		TagSpecifications: []*ec2.TagSpecification{
// 			{
// 				ResourceType: aws.String("instance"),
// 				Tags: []*ec2.Tag{
// 					{
// 						Key:   aws.String("AppID"),
// 						Value: aws.String(appId),
// 					},
// 				},
// 			},
// 		},
// 	})

// 	var params = UpdateEc2StatusParams{
// 		runInstanceResult: runInstanceResult,
// 		keyPairResult:     keyPairResult,
// 		svc:               svc,
// 	}

// 	done := updateEc2Status(params, Db)

// 	if err != nil {
// 		fmt.Println("Couldn't create instance", err)
// 	}
// }

// func updateEc2Status(params UpdateEc2StatusParams, db *gorm.DB) chan bool {
// 	done := make(chan bool)

// 	go func() {

// 		runInstanceResult := params.runInstanceResult
// 		keyPairResult := params.keyPairResult
// 		svc := params.svc
// 		tx := params.tx

// 		// Wait for instance to start
// 		instanceId := runInstanceResult.Instances[0].InstanceId
// 		defaultGroupId := runInstanceResult.Instances[0].SecurityGroups[0].GroupId

// 		fmt.Println("Created instance", *instanceId)

// 		instanceIds := []*string{instanceId}

// 		statusInput := ec2.DescribeInstancesInput{
// 			InstanceIds: instanceIds,
// 		}

// // 3000-6000
// // 8500-8500
// // 8300-8301
// // 443-443
// // 80-80
// // 22-22

// 		// Update ssh inbound group rule
// 		keySecurityGroupsResult, err := svc.UpdateSecurityGroupRuleDescriptionsIngress(&ec2.UpdateSecurityGroupRuleDescriptionsIngressInput{
// 			GroupId: aws.String(*defaultGroupId),
// 			IpPermissions: []*ec2.IpPermission{
// 				{
// 					IpProtocol: aws.String("tcp"),
// 					FromPort:   aws.Int64(22),
// 					ToPort:     aws.Int64(22),
// 					IpRanges: []*ec2.IpRange{
// 						{
// 							CidrIp: aws.String("0.0.0.0/0"),
// 						},
// 					},
// 				},
// 			},
// 		})

// 		if err != nil {
// 			fmt.Println("Could not create keySecurityGroupsResult", err)
// 			return
// 		}
// 		fmt.Println("keySecurityGroupsResult: ", keySecurityGroupsResult)

// 		fmt.Println("waiting for instances to run...")

// 		instanceRunningErr := svc.WaitUntilInstanceRunning(&statusInput)
// 		if instanceRunningErr != nil {
// 			fmt.Printf("failed to wait until instances running: %v", instanceRunningErr)
// 			return
// 		}

// 		fmt.Println("Describing existing instances ...")

// 		description, descriptionErr := svc.DescribeInstances(&statusInput)
// 		if descriptionErr != nil {
// 			fmt.Printf("failed to describe instances: %v", descriptionErr)
// 			return
// 		}

// 		PublicDNSName := description.Reservations[0].Instances[0].PublicDnsName
// 		PublicIPAddress := description.Reservations[0].Instances[0].PublicIpAddress
// 		State := *description.Reservations[0].Instances[0].State

// 		instanceData, err := json.Marshal(description.Reservations[0].Instances[0])

// 		if err != nil {
// 			fmt.Println("=== instanceData json.Marshal error ===", err)
// 			return
// 		}

// 		// Setup host
// 		apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")

// 		api, err := cloudflare.NewWithAPIToken(apiToken)

// 		if err != nil {
// 			log.Fatal(err)
// 			deleteInstanceAndDns(*instanceId, svc, "", nil)
// 			tx.Rollback()
// 			return
// 		}

// 		zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

// 		fmt.Printf("=== PublicDNSName DNSRecord  === %s", *PublicDNSName)
// 		fmt.Printf("=== BackendUrl DNSRecord  === %s", b.BackendUrl)

// 		dnsRecord := cloudflare.DNSRecord{
// 			Type:    "CNAME",
// 			ZoneID:  zoneID,
// 			Name:    b.BackendUrl,
// 			Content: *PublicDNSName,
// 			Proxied: newTrue(),
// 		}

// 		dnsRes, dnsCreateErr := api.CreateDNSRecord(context.Background(), zoneID, dnsRecord)
// 		if dnsCreateErr != nil {
// 			fmt.Println("=== CreateDNSRecord error ===", dnsCreateErr)
// 			deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		dnsRecordData, err := json.Marshal(dnsRes.Result)
// 		if err != nil {
// 			fmt.Println("=== dnsRecordData json.Marshal error ===", err)
// 			deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		dnsInfoResult := tx.Create(&models.DNSRecordInfo{
// 			ID:            nanoid.New(),
// 			AppID:         b.AppID,
// 			EnvironmentID: b.EnvironmentID,
// 			DNSRecordID:   dnsRes.Result.ID,
// 			Content:       dnsRes.Result.Content,
// 			Domain:        dnsRes.Result.Name,
// 			DNSRecord:     dnsRecordData,
// 		})

// 		if dnsInfoResult.Error != nil {
// 			fmt.Printf("Error saving dns instance data == %v", dnsInfoResult.Error)
// 			deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		fmt.Println("=== Created DNSRecord  === ", dnsRes.Result.ID)

// 		VMInfoID := nanoid.New()
// 		VMInfoResult := tx.Create(&models.VMInfo{
// 			ID:                 VMInfoID,
// 			AppID:              b.AppID,
// 			PublicDNSName:      *PublicDNSName,
// 			PublicIPAddress:    *PublicIPAddress,
// 			PublicDNSUser:      "ubuntu",
// 			PemPrivateKeyValue: *keyPairResult.KeyMaterial,
// 			InstanceID:         *instanceId,
// 			InstanceData:       instanceData,
// 			Status:             *State.Code,
// 		})

// 		if VMInfoResult.Error != nil {
// 			fmt.Println("=== VMInfoResult.Error ===", VMInfoResult.Error.Error())
// 			deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		// ================= CONNECT TO EC2 INSTANCE RUN BASH FOR ADD ENV USER ==================

// 		port := "3000"

// 		// password := generatePassword(8, 2, 2, 2)

// 		var connectParams = connect_ec2_instance.Ec2AddEnvUserParams{
// 			PublicDnsNameWithUser: "ubuntu@" + *PublicDNSName,
// 			PemKeyValue:           *keyPairResult.KeyMaterial,
// 			AppID:                 b.AppID,
// 			EnvironmentName:       b.EnvironmentName,
// 			Lang:                  b.Lang,
// 			DomainUrl:             b.BackendUrl,
// 			Port:                  port,
// 			Init:                  "init",
// 			// Password:              password,
// 		}

// 		// Calling Sleep method to wait for last ssh steady

// 		var setupErr error
// 		for i := 1; i <= 3; i++ {
// 			time.Sleep(time.Second * 10)

// 			setupErr = connect_ec2_instance.Ec2AddEnvUser(connectParams)

// 			if checkIsConnectionErr(setupErr) {
// 				fmt.Printf("Error connecting ec2 instance == %v \n", setupErr.Error())
// 				fmt.Printf("==== Try count ====> %v \n", i)
// 			} else {
// 				fmt.Println("== connection success === ")
// 				break // break here
// 			}
// 		}

// 		if checkIsConnectionErr(setupErr) {
// 			fmt.Println("=== ssh connect init setup error Rollback ===", setupErr.Error())
// 			deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		// hashPwd, err := HashPassword(password)
// 		// if err != nil {
// 		// 	fmt.Printf("Error password hash == %v", err)
// 		// 	return
// 		// }

// 		VMUserInfoID := nanoid.New()

// 		VmUserInfoResult := tx.Create(&models.VMUserInfo{
// 			ID:            VMUserInfoID,
// 			EnvironmentID: b.EnvironmentID,
// 			UserName:      b.EnvironmentName,
// 			VMInfoID:      VMInfoID,
// 			// Password:   password,
// 			CurrentlyDeployedFolder: 0,
// 		})

// 		if VmUserInfoResult.Error != nil {
// 			fmt.Printf("Error saving VMInfo instance data == %v", VmUserInfoResult.Error)
// 			return
// 		}

// 		// Calling Sleep method to wait for last ssh steady

// 		var deployOperationRes error
// 		for i := 1; i <= 3; i++ {
// 			time.Sleep(time.Second * 10)
// 			deployOperationRes = vm_operations.CopyBlocksToHandler(vm_operations.RequestObject{
// 				EnvironmentID:  b.EnvironmentID,
// 				AppID:          b.AppID,
// 				DeployId:       b.DeployId,
// 				DeploymentType: b.DeploymentType,
// 				Port:           port,
// 				Tx:             tx,
// 			})

// 			if checkIsConnectionErr(deployOperationRes) {
// 				fmt.Printf("Error deployOperationRes Err == %v", deployOperationRes)
// 				fmt.Printf("==== Try count ====> %v \n", i)
// 			} else {
// 				fmt.Println("== connection success === ")
// 				break // break here
// 			}
// 		}

// 		if checkIsConnectionErr(deployOperationRes) {
// 			fmt.Println("=== ssh deploy operation error and rollback ===")
// 			// deleteInstanceAndDns(*instanceId, svc, dnsRes.Result.ID, api)
// 			tx.Rollback()
// 			return
// 		}

// 		tx.Commit()

// 		fmt.Println("=== End ===")

// 		fmt.Println("=== Health check ===")
// 		time.Sleep(time.Second * 30)
// 		// health_check_operations.Handler(health_check_operations.RequestObject{DeployId: b.DeployId}, db)
// 		fmt.Println("=== Health check End ===")
// 	}()
// 	return done
// }

// func checkIsConnectionErr(err error) bool {
// 	return err != nil && (strings.Contains(err.Error(), "lost connection") || strings.Contains(err.Error(), "connection denied"))
// }
