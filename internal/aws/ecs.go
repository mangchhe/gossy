package aws

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ClusterInfo struct {
	Name string
	Arn  string
}

type TaskInfo struct {
	TaskArn            string
	TaskDefinition     string
	TaskDefinitionName string
	ServiceName        string
	LastStatus         string
	DesiredStatus      string
}

type ContainerInfo struct {
	Name   string
	Status string
}

func GetECSClusters(profile string) ([]ClusterInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	ecsClient := ecs.NewFromConfig(cfg)

	result, err := ecsClient.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	if len(result.ClusterArns) == 0 {
		return []ClusterInfo{}, nil
	}

	describeResult, err := ecsClient.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: result.ClusterArns,
	})
	if err != nil {
		return nil, err
	}

	var clusters []ClusterInfo
	for _, cluster := range describeResult.Clusters {
		clusters = append(clusters, ClusterInfo{
			Name: *cluster.ClusterName,
			Arn:  *cluster.ClusterArn,
		})
	}

	return clusters, nil
}

func GetECSTasks(profile string, clusterArn string) ([]TaskInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	ecsClient := ecs.NewFromConfig(cfg)

	result, err := ecsClient.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster: &clusterArn,
	})
	if err != nil {
		return nil, err
	}

	if len(result.TaskArns) == 0 {
		return []TaskInfo{}, nil
	}

	describeResult, err := ecsClient.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Cluster: &clusterArn,
		Tasks:   result.TaskArns,
	})
	if err != nil {
		return nil, err
	}

	var tasks []TaskInfo
	for _, task := range describeResult.Tasks {
		taskDefName := ExtractTaskDefinitionName(*task.TaskDefinitionArn)
		
		serviceName := "(standalone task)"
		if task.Group != nil && *task.Group != "" {
			if len(*task.Group) > 8 && (*task.Group)[:8] == "service:" {
				serviceName = (*task.Group)[8:]
			} else {
				serviceName = *task.Group
			}
		}
		
		tasks = append(tasks, TaskInfo{
			TaskArn:            *task.TaskArn,
			TaskDefinition:     *task.TaskDefinitionArn,
			TaskDefinitionName: taskDefName,
			ServiceName:        serviceName,
			LastStatus:         *task.LastStatus,
			DesiredStatus:      *task.DesiredStatus,
		})
	}

	return tasks, nil
}

func GetECSContainers(profile string, clusterArn string, taskArn string) ([]ContainerInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	ecsClient := ecs.NewFromConfig(cfg)

	result, err := ecsClient.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Cluster: &clusterArn,
		Tasks:   []string{taskArn},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Tasks) == 0 {
		return []ContainerInfo{}, nil
	}

	var containers []ContainerInfo
	for _, container := range result.Tasks[0].Containers {
		containers = append(containers, ContainerInfo{
			Name:   *container.Name,
			Status: *container.LastStatus,
		})
	}

	return containers, nil
}

func StartECSSession(profile string, clusterName string, taskArn string, containerName string) error {
	cmd := exec.Command("aws", "ecs", "execute-command",
		"--cluster", clusterName,
		"--task", taskArn,
		"--container", containerName,
		"--interactive",
		"--command", "/bin/sh",
		"--profile", profile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Starting ECS session to container %s in task %s...\n", containerName, taskArn)
	return cmd.Run()
}

// ExtractTaskIdFromArn extracts the task ID from a full task ARN
// Example: arn:aws:ecs:region:account:task/cluster-name/1e311516148441978c5f69dd75dbd84a -> 1e311516148441978c5f69dd75dbd84a
func ExtractTaskIdFromArn(taskArn string) string {
	parts := strings.Split(taskArn, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return taskArn
}

// ExtractTaskDefinitionName extracts the task definition name from a full task definition ARN
// Example: arn:aws:ecs:region:account:task-definition/my-app:123 -> my-app
func ExtractTaskDefinitionName(taskDefArn string) string {
	parts := strings.Split(taskDefArn, "/")
	if len(parts) >= 2 {
		taskDefWithRevision := parts[len(parts)-1]
		if colonIndex := strings.LastIndex(taskDefWithRevision, ":"); colonIndex != -1 {
			return taskDefWithRevision[:colonIndex]
		}
		return taskDefWithRevision
	}
	return taskDefArn
}
