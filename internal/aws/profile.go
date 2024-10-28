package aws

import (
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/ini.v1"
)

func GetAWSProfiles() ([]string, error) {
	home, _ := os.UserHomeDir()
	awsCredentialsPath := filepath.Join(home, ".aws", "credentials")

	cfg, err := ini.Load(awsCredentialsPath)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, section := range cfg.Sections() {
		if section.Name() != "DEFAULT" {
			profiles = append(profiles, section.Name())
		}
	}

	sort.Strings(profiles)
	return profiles, nil
}
