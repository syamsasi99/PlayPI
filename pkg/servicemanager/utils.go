package servicemanager

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// IsPortAvailable checks if a port is available for use
func IsPortAvailable(port int) bool {
	timeout := time.Millisecond * 100
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), timeout)
	if err != nil {
		return true // Port is available
	}
	conn.Close()
	return false // Port is in use
}

// KillProcessOnPort kills the process using the specified port
func KillProcessOnPort(port int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		// Use lsof to find and kill the process
		cmd = exec.Command("sh", "-c", fmt.Sprintf("lsof -ti :%d | xargs kill -9 2>/dev/null || true", port))
	case "windows":
		// Use netstat to find the process
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("for /f \"tokens=5\" %%a in ('netstat -aon ^| findstr :%d') do taskkill /F /PID %%a", port))
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// FreePort attempts to free a port by killing the process using it
func FreePort(port int) error {
	if IsPortAvailable(port) {
		return nil // Port is already available
	}

	fmt.Printf("Port %d is in use, attempting to free it...\n", port)
	if err := KillProcessOnPort(port); err != nil {
		return fmt.Errorf("failed to free port %d: %v", port, err)
	}

	// Wait a moment for the port to be freed
	time.Sleep(500 * time.Millisecond)

	if !IsPortAvailable(port) {
		return fmt.Errorf("port %d is still in use after cleanup attempt", port)
	}

	fmt.Printf("✓ Port %d freed successfully\n", port)
	return nil
}

// FreeDashboardPorts frees all ports required by the dashboard
func FreeDashboardPorts() error {
	ports := []int{8000, 8080, 8081, 8082, 8084, 8085, 8086}
	fmt.Println("Checking and freeing dashboard ports...")

	hasErrors := false
	for _, port := range ports {
		if err := FreePort(port); err != nil {
			fmt.Printf("Warning: %v\n", err)
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("some ports could not be freed")
	}

	fmt.Println("✓ All dashboard ports are available")
	return nil
}

// GetProcessOnPort returns the PID of the process using the specified port
func GetProcessOnPort(port int) (int, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	case "windows":
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("for /f \"tokens=5\" %%a in ('netstat -aon ^| findstr :%d') do @echo %%a", port))
	default:
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return 0, fmt.Errorf("no process found on port %d", port)
	}

	pid, err := strconv.Atoi(strings.Split(pidStr, "\n")[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse PID: %v", err)
	}

	return pid, nil
}
