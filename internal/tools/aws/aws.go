package aws

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
)

type awsClient struct {
	ecsClient *ecs.Client
	ssmClient *ssm.Client
	ec2Client *ec2.Client
}

func (awsClient *awsClient) StartBastionProxy(env domain.Environment) {
	result, err := awsClient.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
				},
			},
		},
	})
	if err != nil {
		log.Panicln(err)
	}
	var selectedInstanceId string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.KeyName != nil {
				fmt.Printf("Instance %s: %s\n", *instance.KeyName, *instance.InstanceId)
				if strings.HasPrefix(*instance.KeyName, env.String()) {
					selectedInstanceId = *instance.InstanceId
				}
			} else {
				fmt.Printf("Skipping %s\n", *instance.InstanceId)
			}
		}
	}
	fmt.Printf("Selected instance: %s\n", selectedInstanceId)

	go func() {
		cmd := exec.Command("aws", "ssm", "start-session",
			"--target", selectedInstanceId,
			"--profile", domain.PROFILE,
			"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
			"--parameters", fmt.Sprintf(domain.PARAMETERS),
		)
		fmt.Printf("Starting proxy at port %d\n", domain.SSM_PROXY_PORT)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil {

			fmt.Printf("Command: \n%s\n", cmd.String())
			fmt.Println(err)
			fmt.Println("out:", outb.String(), "err:", errb.String())
		}
	}()

}

func getConfig() aws.Config {
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile(domain.PROFILE),
	)
	return cfg
}

func GetAwsClient() *awsClient {

	awsClient := &awsClient{
		ecsClient: ecs.NewFromConfig(getConfig()),
		ssmClient: ssm.NewFromConfig(getConfig()),
		ec2Client: ec2.NewFromConfig(getConfig()),
	}
	_, err := awsClient.GetEcsClusterMap()
	if err != nil {
		fmt.Println("Not signed it. Triggering login process...")
		cmd := exec.Command("aws", "sso", "login", "--profile", domain.PROFILE)
		cmd.Run()
	}
	return awsClient
}

func (awsClient *awsClient) GetEcsClusterMap() (map[string]string, error) {
	clusters, err := awsClient.ecsClient.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	clusterMap := map[string]string{}
	if err != nil {
		return clusterMap, err
	}

	for _, cluster := range clusters.ClusterArns {
		if strings.Contains(cluster, "dsi-dev") {
			clusterMap["dev"] = cluster
		}
		if strings.Contains(cluster, "dsi-demo") {
			clusterMap["demo"] = cluster
		}
		if strings.Contains(cluster, "dsi-prod") {
			clusterMap["prod"] = cluster
		}
	}
	return clusterMap, nil
}

func (awsClient *awsClient) GetEcsServices(cluster string) map[string]string {
	params := &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(100),
	}
	res, err := awsClient.ecsClient.ListServices(context.TODO(), params)
	if err != nil {
		panic("describe error, " + err.Error())
	}
	serviceMap := map[string]string{}
	for _, service := range res.ServiceArns {
		serviceMap[service[strings.LastIndex(service, "/")+1:]] = service
	}
	return serviceMap
}
