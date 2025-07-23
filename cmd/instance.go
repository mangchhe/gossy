package cmd

import (
	"fmt"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"gossy/internal/util"

	"github.com/AlecAivazis/survey/v2"
)

func chooseInstance(profile string, lastInstanceID string) {
	serviceOptions := []string{util.FormatServiceType("EC2 Instance"), util.FormatServiceType("ECS Pod")}
	var selectedService string
	servicePrompt := &survey.Select{
		Message: "Select service type:",
		Options: serviceOptions,
	}
	err := survey.AskOne(servicePrompt, &selectedService)
	if err != nil {
		util.PrintError("Service selection cancelled")
		return
	}

	if selectedService == util.FormatServiceType("ECS Pod") {

		chooseECSConnection(profile)
		return
	}


	var selectedInstanceID string
	if lastInstanceID != "" {
		selectedInstanceID = lastInstanceID
	} else {
		instances, err := aws.GetEC2Instances(profile)
		if err != nil {
			util.PrintError(fmt.Sprintf("Failed to fetch EC2 instances: %v", err))
			return
		}

		if len(instances) == 0 {
			util.PrintWarning("No EC2 instances found for this profile.")
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
			statusFormatted := util.FormatStatus(instance.State)
			instanceOptions[i] = fmt.Sprintf("%-*s (%s) - %s", maxNameLength, instance.Name, instance.ID, statusFormatted)
		}

		var selectedInstanceIndex int
		instancePrompt := &survey.Select{
			Message: "Select an EC2 Instance:",
			Options: instanceOptions,
		}
		err = survey.AskOne(instancePrompt, &selectedInstanceIndex)
		if err != nil {
			util.PrintError("Instance selection cancelled")
			return
		}

		selectedInstanceID = instances[selectedInstanceIndex].ID
	}


	connectionOptions := []string{util.FormatConnectionType("Instance (SSM)"), util.FormatConnectionType("DB (Port Forwarding)")}
	var connectionType string
	connectionPrompt := &survey.Select{
		Message: "Select connection type:",
		Options: connectionOptions,
	}
	err = survey.AskOne(connectionPrompt, &connectionType)
	if err != nil {
		util.PrintError("Connection type selection cancelled")
		return
	}

	if connectionType == util.FormatConnectionType("Instance (SSM)") {

		err = storage.RecordLastSession(storage.LastSession{
			Profile:    profile,
			InstanceID: selectedInstanceID,
			DatabaseID: "",
		})

		if err != nil {
			util.PrintWarning(fmt.Sprintf("Could not save session: %v", err))
		}


		if err := aws.StartSession(profile, selectedInstanceID); err != nil {
			util.PrintError(fmt.Sprintf("Failed to start SSM session: %v", err))
		} else {
			util.PrintSuccess("SSM session established!")
		}
	} else if connectionType == util.FormatConnectionType("DB (Port Forwarding)") {

		chooseRDSInstance(profile, selectedInstanceID)
	}
}
