package util

import (
	"fmt"
	"log"
	"net"
)

func GetAvailableLocalPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("failed to close listener: %v", err)
		}
	}(listener)

	return listener.Addr().(*net.TCPAddr).Port, nil
}
