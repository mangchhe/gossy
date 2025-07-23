package cmd

import (
	"fmt"
	"strings"
	"github.com/AlecAivazis/survey/v2"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"gossy/internal/util"
)

func Run() {
	util.InitializeScreen()

	profiles, err := aws.GetAWSProfiles()
	if err != nil {
		util.PrintError(fmt.Sprintf("Failed to fetch AWS profiles: %v", err))
		return
	}

	if len(profiles) == 0 {
		util.PrintWarning("No AWS profiles found.")
		util.PrintWarning("Run: aws configure --profile <profile-name>")
		return
	}

	recentProfiles := getRecentProfiles(profiles)
	profileOptions := append([]string{"Connect with Last Session"}, recentProfiles...)

	var selectedProfile string
	profilePrompt := &survey.Select{
		Message:  "Select an AWS Profile:",
		Options:  profileOptions,
		PageSize: len(profileOptions),
	}
	err = survey.AskOne(profilePrompt, &selectedProfile)
	if err != nil {
		util.PrintError("Profile selection cancelled")
		return
	}

	if selectedProfile != "Connect with Last Session" {
		util.UpdateProfile(selectedProfile)
		util.UpdateCommand("aws ssm start-session --profile " + selectedProfile)
	}

	if selectedProfile == "Connect with Last Session" {
		lastSession, err := storage.LoadLastSession()
		if err == nil {
			util.UpdateProfile(lastSession.Profile)
			util.UpdateCommand("aws ssm start-session --profile " + lastSession.Profile)
			if lastSession.DatabaseID != "" {
				rdsInstance, err := aws.GetRDSInstanceByDatabaseID(lastSession.Profile, lastSession.DatabaseID)
				if err == nil {
					localPort, _ := util.GetAvailableLocalPort()
					command := fmt.Sprintf("aws ssm start-session --profile %s --target %s --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{\"portNumber\":[\"%d\"],\"localPortNumber\":[\"%d\"],\"host\":[\"%s\"]}'", 
						lastSession.Profile, lastSession.InstanceID, rdsInstance.Port, localPort, rdsInstance.Endpoint)
					util.UpdateCommand(command)
					util.ClearContent()
					util.Command.Printf("Command: %s\n", command)
					fmt.Println(strings.Repeat("─", 80))
					fmt.Println()
				}
				err = aws.StartRDSConnection(lastSession.Profile, lastSession.InstanceID, lastSession.DatabaseID)
				if err != nil {
					util.PrintError(fmt.Sprintf("Failed to start DB session: %v", err))
				} else {
					util.PrintSuccess("RDS connection established!")
				}
			} else if lastSession.InstanceID != "" {
				command := fmt.Sprintf("aws ssm start-session --target %s --profile %s", lastSession.InstanceID, lastSession.Profile)
				util.UpdateCommand(command)
				util.ClearContent()
				util.Command.Printf("Command: %s\n", command)
				fmt.Println(strings.Repeat("─", 80))
				fmt.Println()
				err = aws.StartInstanceConnection(lastSession.Profile, lastSession.InstanceID)
				if err != nil {
					util.PrintError(fmt.Sprintf("Failed to start EC2 session: %v", err))
				} else {
					util.PrintSuccess("EC2 connection established!")
				}
			}
			return
		}
		util.ClearContent()
		util.PrintWarning("No previous session found. Starting fresh...")
		fmt.Println()
	}
	err = storage.RecordLastSession(storage.LastSession{
		Profile:    selectedProfile,
		InstanceID: "",
		DatabaseID: "",
	})
	if err != nil {
		util.PrintWarning(fmt.Sprintf("Could not save session: %v", err))
	}

	chooseInstance(selectedProfile, "")
}
