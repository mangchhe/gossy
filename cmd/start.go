package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"gossy/internal/util"
)

func Run() {
	util.PrintWelcome()
	

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

	if selectedProfile == "Connect with Last Session" {

		lastSession, err := storage.LoadLastSession()
		if err == nil {


			if lastSession.DatabaseID != "" {
				err = aws.StartRDSConnection(lastSession.Profile, lastSession.InstanceID, lastSession.DatabaseID)
				if err != nil {
					util.PrintError(fmt.Sprintf("Failed to start DB session: %v", err))
				} else {
					util.PrintSuccess("RDS connection established!")
				}
			} else if lastSession.InstanceID != "" {
				err = aws.StartInstanceConnection(lastSession.Profile, lastSession.InstanceID)
				if err != nil {
					util.PrintError(fmt.Sprintf("Failed to start EC2 session: %v", err))
				} else {
					util.PrintSuccess("EC2 connection established!")
				}
			}
			return
		}
		util.PrintWarning("No previous session found. Starting fresh...")
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
