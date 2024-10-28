package aws

import (
	"fmt"
	"gossy/internal/storage"
	"gossy/internal/util"
)

func StartInstanceConnection(profile, instanceID string) error {
	fmt.Printf("Starting SSM session to instance %s...\n", instanceID)
	if err := StartSession(profile, instanceID); err != nil {
		return fmt.Errorf("Failed to start SSM session: %v", err)
	}
	return nil
}

func StartRDSConnection(profile, instanceID, databaseID string) error {
	selectedRDS, err := GetRDSInstanceByDatabaseID(profile, databaseID)
	if err != nil {
		return fmt.Errorf("RDS instance with ID %s not found: %v", databaseID, err)
	}

	localPort, err := util.GetAvailableLocalPort()
	if err != nil {
		return fmt.Errorf("Error finding available local port: %v", err)
	}

	err = storage.RecordLastSession(storage.LastSession{
		Profile:    profile,
		InstanceID: instanceID,
		DatabaseID: selectedRDS.ID,
	})

	if err != nil {
		fmt.Printf("Failed to record last session: %v\n", err)
	}

	if err := StartPortForwardingSession(profile, instanceID, selectedRDS.Endpoint, selectedRDS.Port, int32(localPort)); err != nil {
		return fmt.Errorf("Failed to start port-forwarding session: %v", err)
	}

	fmt.Printf("Port forwarding started on local port %d for RDS instance %s\n", localPort, selectedRDS.ID)
	return nil
}
