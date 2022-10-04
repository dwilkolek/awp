package aws

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type awsClient struct {
	*ecs.Client
}

func GetAwsClient() *awsClient {
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile("developer"),
	)

	awsClient := &awsClient{
		ecs.NewFromConfig(cfg),
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
	clusters, err := awsClient.ListClusters(context.TODO(), &ecs.ListClustersInput{})
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
	res, err := awsClient.ListServices(context.TODO(), params)
	if err != nil {
		panic("describe error, " + err.Error())
	}
	serviceMap := map[string]string{}
	for _, service := range res.ServiceArns {
		serviceMap[service[strings.LastIndex(service, "/")+1:]] = service
	}
	return serviceMap
}
