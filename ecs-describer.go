package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"ecs-describer/utils"

	"github.com/olekukonko/tablewriter"
)

var clusterName = flag.String("cluster", "", "ECS cluster name")

func main() {
	flag.Parse()
	if *clusterName == "" {
		fmt.Fprintf(os.Stderr, "Error: Must specify cluster name")
		os.Exit(1)
	}

	fmt.Printf("Retrieving cluster details...\n")
	start := time.Now()

	taskDetails := utils.DescribeClusterTasks(clusterName)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Service", "Container", "Image", "Instance Id", "IP"})

	for _, taskDetail := range taskDetails.Tasks {
		for _, containerDetails := range taskDetail.Containers {

			ContainerArn := []string{*taskDetail.ContainerInstanceArn}
			instanceDetails := utils.DescribeClusterInstances(clusterName, ContainerArn)

			row := []string{*taskDetail.Group, *containerDetails.Name, *containerDetails.Image,
				*instanceDetails.Reservations[0].Instances[0].InstanceId,
				*instanceDetails.Reservations[0].Instances[0].PublicIpAddress}
			table.Append(row)
		}
	}
	table.Render()
	duration := time.Since(start)
	fmt.Printf("\nCluster details retrieved in %v", duration)
}
