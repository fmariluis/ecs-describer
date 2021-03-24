package utils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var ecsClient = getECSClient()
var ec2Client = getEC2Client()

func getAWSConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Can't config AWS")
	}
	return cfg
}

func getECSClient() ecs.Client {
	cfg := getAWSConfig()
	ecsClient := ecs.NewFromConfig(cfg)
	return *ecsClient
}

func getEC2Client() ec2.Client {
	cfg := getAWSConfig()
	ec2Client := ec2.NewFromConfig(cfg)
	return *ec2Client
}

// Returns a list of running tasks in a given cluster
func DescribeClusterTasks(clusterName *string) ecs.DescribeTasksOutput {
	tasks, err := ecsClient.ListTasks(context.TODO(),
		&ecs.ListTasksInput{Cluster: clusterName, MaxResults: aws.Int32(20), DesiredStatus: "RUNNING"})
	if err != nil {
		log.Fatalf("failed to list tasks, %v", err)
	}

	taskDetails, err := ecsClient.DescribeTasks(context.TODO(),
		&ecs.DescribeTasksInput{Tasks: tasks.TaskArns, Cluster: clusterName})

	if err != nil {
		log.Fatalf("failed to retrieve task details, %v", err)
	}

	return *taskDetails
}

// Returns details about the instances hosting the ECS services
func DescribeClusterInstances(clusterName *string, ContainerArn []string) ec2.DescribeInstancesOutput {
	containerInstances, err := ecsClient.DescribeContainerInstances(context.TODO(),
		&ecs.DescribeContainerInstancesInput{Cluster: clusterName, ContainerInstances: ContainerArn})
	if err != nil {
		log.Fatalf("failed to retrieve container instance details, %v", err)
	}
	instanceID := *containerInstances.ContainerInstances[0].Ec2InstanceId

	instanceIDs := []string{instanceID}
	instanceDetails, err := ec2Client.DescribeInstances(context.TODO(),
		&ec2.DescribeInstancesInput{InstanceIds: instanceIDs})

	if err != nil {
		log.Fatalf("failed to retrieve container instance details, %v", err)
	}

	return *instanceDetails
}
