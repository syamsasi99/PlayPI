package servicemanager

import (
	"fmt"
	"net"
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
