package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"gossy/internal/util"
)

func chooseRDSInstance(profile string, instanceID string) {
	rdsInstances, err := aws.GetRDSInstancesByProfile(profile)
	if err != nil {
		fmt.Printf("Error fetching RDS instances: %v\n", err)
		return
	}

	if len(rdsInstances) == 0 {
		fmt.Println("No RDS instances found.")
		return
	}

	rdsOptions := []string{}
	for _, rds := range rdsInstances {
		rdsOptions = append(rdsOptions, fmt.Sprintf("%s (%s:%d)", rds.ID, rds.Endpoint, rds.Port))
	}

	var selectedRDSIndex int
	rdsPrompt := &survey.Select{
		Message: "Select an RDS instance for port forwarding:",
		Options: rdsOptions,
	}
	err = survey.AskOne(rdsPrompt, &selectedRDSIndex)
	if err != nil {
		fmt.Printf("Error selecting RDS instance: %v\n", err)
		return
	}

	selectedRDS := rdsInstances[selectedRDSIndex]
	localPort := selectedRDS.Port

	parameters := fmt.Sprintf(`'{"portNumber":["%d"],"localPortNumber":["%d"],"host":["%s"]}'`, selectedRDS.Port, localPort, selectedRDS.Endpoint)
	command := fmt.Sprintf("aws ssm start-session --profile %s --target %s --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters %s", profile, instanceID, parameters)
	util.UpdateCommand(command)
	util.PrintFixedCommand("aws ssm start-session", "--profile", profile, "--target", instanceID, "--document-name", "AWS-StartPortForwardingSessionToRemoteHost", "--parameters", parameters)

	err = storage.RecordLastSession(storage.LastSession{
		Profile:    profile,
		InstanceID: instanceID,
		DatabaseID: selectedRDS.ID,
	})

	if err != nil {
	}

	if err := aws.StartPortForwardingSession(profile, instanceID, selectedRDS.Endpoint, selectedRDS.Port, localPort); err != nil {
		fmt.Printf("Failed to start port-forwarding session: %v\n", err)
		return
	}
}
