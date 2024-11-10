package cmd

import (
	"fmt"
	"gossy/internal/aws"
	"gossy/internal/storage"

	"github.com/AlecAivazis/survey/v2"
)

func chooseInstance(profile string, lastInstanceID string) {
	var selectedInstanceID string
	if lastInstanceID != "" {
		selectedInstanceID = lastInstanceID
	} else {
		instances, err := aws.GetEC2Instances(profile)
		if err != nil {
			fmt.Printf("Error fetching EC2 instances: %v\n", err)
			return
		}

		if len(instances) == 0 {
			fmt.Println("No instances found for this profile.")
			return
		}

		maxNameLength := 0
		for _, instance := range instances {
			if len(instance.Name) > maxNameLength {
				maxNameLength = len(instance.Name)
			}
		}

		instanceOptions := make([]string, len(instances))
		for i, instance := range instances {
			instanceOptions[i] = fmt.Sprintf("%-*s (%s) - %s", maxNameLength, instance.Name, instance.ID, instance.State)
		}

		var selectedInstanceIndex int
		instancePrompt := &survey.Select{
			Message: "Select an EC2 Instance:",
			Options: instanceOptions,
		}
		err = survey.AskOne(instancePrompt, &selectedInstanceIndex)
		if err != nil {
			fmt.Println("Error selecting instance.")
			return
		}

		selectedInstanceID = instances[selectedInstanceIndex].ID
	}

	connectionOptions := []string{"Instance (SSM)", "DB (Port Forwarding)"}
	var connectionType string
	connectionPrompt := &survey.Select{
		Message: "Select connection type:",
		Options: connectionOptions,
	}
	err := survey.AskOne(connectionPrompt, &connectionType)
	if err != nil {
		fmt.Println("Error selecting connection type.")
		return
	}

	if connectionType == "Instance (SSM)" {
		err = storage.RecordLastSession(storage.LastSession{
			Profile:    profile,
			InstanceID: selectedInstanceID,
			DatabaseID: "",
		})

		if err != nil {
			fmt.Printf("Failed to record last session: %v\n", err)
		}

		if err := aws.StartSession(profile, selectedInstanceID); err != nil {
			fmt.Printf("Failed to start SSM session: %v\n", err)
		}
	} else if connectionType == "DB (Port Forwarding)" {
		chooseRDSInstance(profile, selectedInstanceID)
	}
}
