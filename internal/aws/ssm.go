package aws

import (
	"fmt"
	"os"
	"os/exec"
)

func StartSession(profile string, instanceID string) error {
	cmd := exec.Command("aws", "ssm", "start-session",
		"--target", instanceID,
		"--profile", profile)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Starting SSM session to instance %s...\n", instanceID)
	return cmd.Run()
}

func StartPortForwardingSession(profile string, instanceID string, remoteHost string, remotePort, localPort int32) error {
	cmd := exec.Command("aws", "ssm", "start-session",
		"--profile", profile,
		"--target", instanceID,
		"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters", fmt.Sprintf(`{"portNumber":["%d"],"localPortNumber":["%d"],"host":["%s"]}`, remotePort, localPort, remoteHost),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Starting port-forwarding session to RDS instance at %s:%d via instance %s (local port %d)...\n", remoteHost, remotePort, instanceID, localPort)
	return cmd.Run()
}
