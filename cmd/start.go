package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gossy/internal/aws"
	"gossy/internal/storage"
)

func Run() {
	profiles, err := aws.GetAWSProfiles()
	if err != nil {
		fmt.Printf("Error fetching profiles: %v\n", err)
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
		fmt.Printf("Error selecting profile: %v\n", err)
		return
	}

	if selectedProfile == "Connect with Last Session" {
		lastSession, err := storage.LoadLastSession()
		if err == nil {
			fmt.Printf("Connecting to last session: Profile=%s, InstanceID=%s, DatabaseID=%s\n", lastSession.Profile, lastSession.InstanceID, lastSession.DatabaseID)

			if lastSession.DatabaseID != "" {
				err = aws.StartRDSConnection(lastSession.Profile, lastSession.InstanceID, lastSession.DatabaseID)
				if err != nil {
					fmt.Printf("Failed to start DB port-forwarding session: %v\n", err)
				}
			} else if lastSession.InstanceID != "" {
				err = aws.StartInstanceConnection(lastSession.Profile, lastSession.InstanceID)
				if err != nil {
					fmt.Printf("Failed to start EC2 instance session: %v\n", err)
				}
			}
			return
		}
		fmt.Println("No previous session found.")
		return
	}

	err = storage.RecordLastSession(storage.LastSession{
		Profile:    selectedProfile,
		InstanceID: "",
		DatabaseID: "",
	})
	if err != nil {
		fmt.Printf("Failed to record last session: %v\n", err)
	}

	chooseInstance(selectedProfile, "")
}
