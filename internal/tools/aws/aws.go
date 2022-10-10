package aws

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type awsClient struct {
	ecsClient *ecs.Client
	ssmClient *ssm.Client
}

func (awsClient *awsClient) StartBastionProxy() {
	so, err := awsClient.ssmClient.StartSession(context.TODO(), &ssm.StartSessionInput{
		Target:       aws.String("i-09e46925fbb91c0f4"),
		DocumentName: aws.String("AWS-StartPortForwardingSession"),
		Parameters: map[string][]string{
			"portNumber":      {"80"},
			"localPortNumber": {"17777"},
		},
	})
	if err != nil {
		log.Panicln(err)
	}

	log.Println(so.SessionId)
}

func getConfig() aws.Config {
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile("developer"),
	)
	return cfg
}

func GetAwsClient() *awsClient {

	awsClient := &awsClient{
		ecsClient: ecs.NewFromConfig(getConfig()),
		ssmClient: ssm.NewFromConfig(getConfig()),
	}
	_, err := awsClient.GetEcsClusterMap()
	if err != nil {
		fmt.Println("Not signed it. Triggering login process...")
		cmd := exec.Command("aws", "sso", "login", "--profile", "developer")
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
