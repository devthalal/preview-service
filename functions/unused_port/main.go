package unused_port

import (
	common "ab-preview-service/common"
	"fmt"
	"net"
)

func GetUnusedPort(portCount int) common.FunctionReturn {
	numUnusedPorts := portCount // Number of unused ports you want to find

	if portCount < 1 {
		numUnusedPorts = 1
	}

	unusedPorts := findNUnusedPorts(numUnusedPorts)

	return common.FunctionReturn{
		Message: "Success",
		Data:    unusedPorts,
	}
}

func findNUnusedPorts(n int) []int {
	unusedPorts := []int{}

	port := 3000 // Start checking from this port number

	for len(unusedPorts) < n {
		if isPortUnused(port) {
			unusedPorts = append(unusedPorts, port)
		}
		port++
	}

	return unusedPorts
}

func isPortUnused(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false // Port is already in use
	}
	defer listener.Close()
	return true // Port is unused
}
