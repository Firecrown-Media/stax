package system

import (
	"testing"
	"time"
)

func TestIsPortAvailable(t *testing.T) {
	// Port 0 should always be available (system assigns a random port)
	// High ports should generally be available
	testPort := 59999

	available := IsPortAvailable(testPort)
	if !available {
		t.Logf("Port %d is not available (this is not necessarily an error)", testPort)
	}
}

func TestFindAvailablePort(t *testing.T) {
	port, err := FindAvailablePort(50000)
	if err != nil {
		t.Fatalf("FindAvailablePort() error = %v", err)
	}

	if port < 50000 {
		t.Errorf("FindAvailablePort(50000) = %d, expected >= 50000", port)
	}

	// Verify the port is actually available
	if !IsPortAvailable(port) {
		t.Errorf("FindAvailablePort returned port %d which is not available", port)
	}
}

func TestGetPortsInUse(t *testing.T) {
	ports := []int{80, 443, 50001, 50002}
	inUse := GetPortsInUse(ports)

	t.Logf("Ports in use: %v", inUse)

	// We can't make assumptions about which ports are in use,
	// but we can verify the function returns a valid slice
	if inUse == nil {
		t.Error("GetPortsInUse() returned nil")
	}
}

func TestWaitForPort(t *testing.T) {
	// Find an available port
	port, err := FindAvailablePort(50000)
	if err != nil {
		t.Fatalf("FindAvailablePort() error = %v", err)
	}

	// Wait for it to be available (it already is)
	err = WaitForPort(port, true, 1*time.Second)
	if err != nil {
		t.Errorf("WaitForPort() error = %v", err)
	}
}

func TestDefaultDDEVPorts(t *testing.T) {
	ports := DefaultDDEVPorts()

	if len(ports) == 0 {
		t.Error("DefaultDDEVPorts() returned empty slice")
	}

	// Verify it contains expected ports
	expectedPorts := map[int]bool{
		80:   true,
		443:  true,
		3306: true,
	}

	for expectedPort := range expectedPorts {
		found := false
		for _, port := range ports {
			if port == expectedPort {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DefaultDDEVPorts() missing expected port %d", expectedPort)
		}
	}
}

func TestRecommendedPorts(t *testing.T) {
	defaults := []int{80, 443}
	recommendations := RecommendedPorts(defaults)

	t.Logf("Recommendations: %v", recommendations)

	// Verify recommendations are valid
	for original, recommended := range recommendations {
		if recommended <= original {
			t.Errorf("Recommended port %d is not greater than original %d", recommended, original)
		}

		if !IsPortAvailable(recommended) {
			t.Errorf("Recommended port %d is not available", recommended)
		}
	}
}
