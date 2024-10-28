package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LastSession struct {
	Profile    string `json:"profile"`
	InstanceID string `json:"instance_id"`
	DatabaseID string `json:"database_id"`
}

var sessionFilePath = filepath.Join(getGossyDir(), "last_session.json")

func getGossyDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to get the user's home directory: " + err.Error())
	}
	gossyDir := filepath.Join(homeDir, ".gossy")

	if _, err := os.Stat(gossyDir); os.IsNotExist(err) {
		if err := os.Mkdir(gossyDir, 0755); err != nil {
			panic("Unable to create .gossy directory: " + err.Error())
		}
	}
	return gossyDir
}

func RecordLastSession(session LastSession) error {
	file, err := os.Create(sessionFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Error closing file: %v\n", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	return encoder.Encode(session)
}

func LoadLastSession() (LastSession, error) {
	var session LastSession
	file, err := os.Open(sessionFilePath)
	if err != nil {
		return session, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Error closing file: %v\n", cerr)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&session)
	return session, err
}
