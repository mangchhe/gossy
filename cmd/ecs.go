package cmd

import (
	"fmt"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"gossy/internal/util"

	"github.com/AlecAivazis/survey/v2"
)

func chooseECSConnection(profile string) {
	// Get ECS clusters

	clusters, err := aws.GetECSClusters(profile)
	if err != nil {
		util.PrintError(fmt.Sprintf("Failed to fetch ECS clusters: %v", err))
		return
	}

	if len(clusters) == 0 {
		util.PrintWarning("No ECS clusters found for this profile.")
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
		util.PrintError("Cluster selection cancelled")
		return
	}

	selectedCluster := clusters[selectedClusterIndex]

	// Get tasks for selected cluster

	tasks, err := aws.GetECSTasks(profile, selectedCluster.Arn)
	if err != nil {
		util.PrintError(fmt.Sprintf("Failed to fetch ECS tasks: %v", err))
		return
	}

	if len(tasks) == 0 {
		util.PrintWarning("No running ECS tasks found in this cluster.")
		return
	}

	// Select task
	taskOptions := make([]string, len(tasks))
	for i, task := range tasks {
		// Extract task ID from ARN for display
		taskId := aws.ExtractTaskIdFromArn(task.TaskArn)
		statusFormatted := util.FormatStatus(task.LastStatus)
		// Display format: "service-name/task-definition-name (task-id) - status"
		taskOptions[i] = fmt.Sprintf("%s/%s (%s) - %s", task.ServiceName, task.TaskDefinitionName, taskId[:8], statusFormatted)
	}

	var selectedTaskIndex int
	taskPrompt := &survey.Select{
		Message: "Select a Task:",
		Options: taskOptions,
	}
	err = survey.AskOne(taskPrompt, &selectedTaskIndex)
	if err != nil {
		util.PrintError("Task selection cancelled")
		return
	}

	selectedTask := tasks[selectedTaskIndex]

	// Get containers for selected task

	containers, err := aws.GetECSContainers(profile, selectedCluster.Arn, selectedTask.TaskArn)
	if err != nil {
		util.PrintError(fmt.Sprintf("Failed to fetch ECS containers: %v", err))
		return
	}

	if len(containers) == 0 {
		util.PrintWarning("No containers found in this task.")
		return
	}

	// Select container
	containerOptions := make([]string, len(containers))
	for i, container := range containers {
		statusFormatted := util.FormatStatus(container.Status)
		containerOptions[i] = fmt.Sprintf("%s - %s", container.Name, statusFormatted)
	}

	var selectedContainerIndex int
	containerPrompt := &survey.Select{
		Message: "Select a Container:",
		Options: containerOptions,
	}
	err = survey.AskOne(containerPrompt, &selectedContainerIndex)
	if err != nil {
		util.PrintError("Container selection cancelled")
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
		util.PrintWarning(fmt.Sprintf("Could not save session: %v", err))
	}

	// Start ECS session
	taskId := aws.ExtractTaskIdFromArn(selectedTask.TaskArn)

	if err := aws.StartECSSession(profile, selectedCluster.Name, taskId, selectedContainer.Name); err != nil {
		util.PrintError(fmt.Sprintf("Failed to start ECS session: %v", err))
	} else {
		util.PrintSuccess("ECS session established!")
	}
}


