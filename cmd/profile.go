package cmd

import (
	"fmt"
	"gossy/internal/aws"
	"gossy/internal/storage"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Select an AWS profile and choose an EC2 instance",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := aws.GetAWSProfiles()
		if err != nil {
			fmt.Printf("Error fetching profiles: %v\n", err)
			return
		}

		var selectedProfile string
		profilePrompt := &survey.Select{
			Message: "Select an AWS Profile:",
			Options: profiles,
		}
		err = survey.AskOne(profilePrompt, &selectedProfile)
		if err != nil {
			fmt.Printf("Error selecting profile: %v\n", err)
			return
		}

		chooseInstance(selectedProfile, "")
	},
}

func getRecentProfiles(profiles []string) []string {
	lastSession, err := storage.LoadLastSession()
	if err != nil {
		return profiles
	}

	return sortProfiles(profiles, lastSession.Profile)
}

func sortProfiles(profiles []string, recentProfile string) []string {
	profileSet := make(map[string]bool)
	for _, profile := range profiles {
		profileSet[profile] = true
	}

	sortedProfiles := []string{}
	if recentProfile != "" && profileSet[recentProfile] {
		sortedProfiles = append(sortedProfiles, recentProfile)
		delete(profileSet, recentProfile)
	}

	remainingProfiles := []string{}
	for profile := range profileSet {
		remainingProfiles = append(remainingProfiles, profile)
	}
	sort.Strings(remainingProfiles)

	return append(sortedProfiles, remainingProfiles...)
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
