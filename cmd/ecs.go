package cmd

import (
	"fmt"
	"gossy/internal/aws"
	"gossy/internal/storage"

	"github.com/AlecAivazis/survey/v2"
)

func chooseECSConnection(profile string) {
	// Get ECS clusters
	clusters, err := aws.GetECSClusters(profile)
	if err != nil {
		fmt.Printf("Error fetching ECS clusters: %v\n", err)
		return
	}

	if len(clusters) == 0 {
		fmt.Println("No ECS clusters found for this profile.")
		return
	}

	// Select cluster
	clusterOptions := make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterOptions[i] = cluster.Name
	}

	var selectedClusterIndex int
	clusterPrompt := &survey.Select{
		Message: "Select an ECS Cluster:",
		Options: clusterOptions,
	}
	err = survey.AskOne(clusterPrompt, &selectedClusterIndex)
	if err != nil {
		fmt.Println("Error selecting cluster.")
		return
	}

	selectedCluster := clusters[selectedClusterIndex]

	// Get tasks for selected cluster
	tasks, err := aws.GetECSTasks(profile, selectedCluster.Arn)
	if err != nil {
		fmt.Printf("Error fetching ECS tasks: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No running tasks found in this cluster.")
		return
	}

	// Select task
	taskOptions := make([]string, len(tasks))
	for i, task := range tasks {
		// Extract task ID from ARN for display
		taskId := aws.ExtractTaskIdFromArn(task.TaskArn)
		// Display format: "service-name/task-definition-name (task-id) - status"
		taskOptions[i] = fmt.Sprintf("%s/%s (%s) - %s", task.ServiceName, task.TaskDefinitionName, taskId[:8], task.LastStatus)
	}

	var selectedTaskIndex int
	taskPrompt := &survey.Select{
		Message: "Select a Task:",
		Options: taskOptions,
	}
	err = survey.AskOne(taskPrompt, &selectedTaskIndex)
	if err != nil {
		fmt.Println("Error selecting task.")
		return
	}

	selectedTask := tasks[selectedTaskIndex]

	// Get containers for selected task
	containers, err := aws.GetECSContainers(profile, selectedCluster.Arn, selectedTask.TaskArn)
	if err != nil {
		fmt.Printf("Error fetching ECS containers: %v\n", err)
		return
	}

	if len(containers) == 0 {
		fmt.Println("No containers found in this task.")
		return
	}

	// Select container
	containerOptions := make([]string, len(containers))
	for i, container := range containers {
		containerOptions[i] = fmt.Sprintf("%s (%s)", container.Name, container.Status)
	}

	var selectedContainerIndex int
	containerPrompt := &survey.Select{
		Message: "Select a Container:",
		Options: containerOptions,
	}
	err = survey.AskOne(containerPrompt, &selectedContainerIndex)
	if err != nil {
		fmt.Println("Error selecting container.")
		return
	}

	selectedContainer := containers[selectedContainerIndex]

	// Record session (we'll store cluster name in InstanceID field and container name in DatabaseID field for ECS sessions)
	err = storage.RecordLastSession(storage.LastSession{
		Profile:    profile,
		InstanceID: selectedCluster.Name,
		DatabaseID: selectedContainer.Name,
	})
	if err != nil {
		fmt.Printf("Failed to record last session: %v\n", err)
	}

	// Start ECS session
	taskId := aws.ExtractTaskIdFromArn(selectedTask.TaskArn)
	if err := aws.StartECSSession(profile, selectedCluster.Name, taskId, selectedContainer.Name); err != nil {
		fmt.Printf("Failed to start ECS session: %v\n", err)
	}
}


