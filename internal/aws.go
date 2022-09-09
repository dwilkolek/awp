package awsserviceproxy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/txn2/txeh"
)

type ShortRunConfiguration struct {
	Env      string
	Services []ServiceConfiguration
}

func SetupHosts() {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	client := getEcsClient()
	clusters, _ := getEcsClusterMap(client)
	services := getEcsServices(client, clusters["dev"])

	knownHosts := []string{}

	hosts.AddHost("127.0.0.1", "awp")

	for service := range services {
		hosts.AddHost("127.0.0.1", service+".service")
		knownHosts = append(knownHosts, service+".service")
	}

	AWPConfig.Hosts = knownHosts
	saveAWPConfig()
	// hfData := hosts.RenderHostsFile()
	// fmt.Println(hfData)
	err = hosts.Save()
	if err != nil {
		panic(err)
	}
}

func getEcsClient() *ecs.Client {
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile("developer"),
	)

	client := ecs.NewFromConfig(cfg)

	_, err := getEcsClusterMap(client)
	if err != nil {
		fmt.Println("Not signed it. Triggering login process...")
		cmd := exec.Command("aws", "sso", "login", "--profile", "developer")
		cmd.Run()
	}
	return client
}

func getEcsClusterMap(client *ecs.Client) (map[string]string, error) {
	clusters, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
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

func getEcsServices(client *ecs.Client, cluster string) map[string]string {
	params := &ecs.ListServicesInput{
		Cluster:    aws.String(cluster),
		MaxResults: aws.Int32(100),
	}
	res, err := client.ListServices(context.TODO(), params)
	if err != nil {
		panic("describe error, " + err.Error())
	}
	serviceMap := map[string]string{}
	for _, service := range res.ServiceArns {
		serviceMap[service[strings.LastIndex(service, "/")+1:]] = service
	}
	return serviceMap
}
