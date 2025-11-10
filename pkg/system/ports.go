package system

import (
	"fmt"
	"net"
	"time"
)

// IsPortAvailable checks if a port is available on localhost
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// FindAvailablePort finds an available port starting from the given port
func FindAvailablePort(startPort int) (int, error) {
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found in range %d-%d", startPort, startPort+maxAttempts)
}

// GetPortsInUse returns a list of ports currently in use
func GetPortsInUse(ports []int) []int {
	var inUse []int
	for _, port := range ports {
		if !IsPortAvailable(port) {
			inUse = append(inUse, port)
		}
	}
	return inUse
}

// WaitForPort waits for a port to become available or occupied
func WaitForPort(port int, available bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		isAvail := IsPortAvailable(port)
		if isAvail == available {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	if available {
		return fmt.Errorf("port %d did not become available within %v", port, timeout)
	}
	return fmt.Errorf("port %d did not become occupied within %v", port, timeout)
}

// CheckRequiredPorts checks if required ports are available
func CheckRequiredPorts(ports []int) ([]int, error) {
	inUse := GetPortsInUse(ports)
	if len(inUse) > 0 {
		return inUse, fmt.Errorf("the following required ports are in use: %v", inUse)
	}
	return nil, nil
}

// GetProcessUsingPort returns information about the process using a port (Unix-like systems)
func GetProcessUsingPort(port int) (string, error) {
	// This is a placeholder - actual implementation would use lsof on Unix
	// or netstat on Windows
	if IsPortAvailable(port) {
		return "", fmt.Errorf("port %d is not in use", port)
	}
	return fmt.Sprintf("Unknown process using port %d", port), nil
}

// DefaultDDEVPorts returns the default ports used by DDEV
func DefaultDDEVPorts() []int {
	return []int{
		80,    // HTTP
		443,   // HTTPS
		3306,  // MySQL
		8025,  // Mailhog
		8036,  // PHPMyAdmin
	}
}

// RecommendedPorts recommends alternative ports if defaults are in use
func RecommendedPorts(defaults []int) map[int]int {
	recommendations := make(map[int]int)
	for _, port := range defaults {
		if !IsPortAvailable(port) {
			// Find an alternative port
			alt, err := FindAvailablePort(port + 1000)
			if err == nil {
				recommendations[port] = alt
			}
		}
	}
	return recommendations
}
