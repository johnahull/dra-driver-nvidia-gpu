package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

)

// getNUMANodeByPCIBusID reads the NUMA node for a PCI device from sysfs.
func getNUMANodeByPCIBusID(pciBusID string) (int, error) {
	path := fmt.Sprintf("/sys/bus/pci/devices/%s/numa_node", pciBusID)
	data, err := os.ReadFile(path)
	if err != nil {
		return -1, fmt.Errorf("reading numa_node for %s: %w", pciBusID, err)
	}
	numa, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return -1, fmt.Errorf("parsing numa_node for %s: %w", pciBusID, err)
	}
	return numa, nil
}

// getSocketByNUMANode reads the CPU socket (physical_package_id) for a given NUMA node.
// Path: numa_node → cpulist → first CPU → physical_package_id
func getSocketByNUMANode(numaNode int) (int, error) {
	cpulistPath := fmt.Sprintf("/sys/devices/system/node/node%d/cpulist", numaNode)
	data, err := os.ReadFile(cpulistPath)
	if err != nil {
		return -1, fmt.Errorf("reading cpulist for NUMA %d: %w", numaNode, err)
	}

	cpulist := strings.TrimSpace(string(data))
	if cpulist == "" {
		return -1, fmt.Errorf("empty cpulist for NUMA %d", numaNode)
	}

	// Parse first CPU from the list (e.g., "0,2,4,6" → "0" or "0-31" → "0")
	firstCPU := cpulist
	if idx := strings.IndexAny(cpulist, ",-"); idx > 0 {
		firstCPU = cpulist[:idx]
	}

	socketPath := fmt.Sprintf("/sys/devices/system/cpu/cpu%s/topology/physical_package_id", firstCPU)
	socketData, err := os.ReadFile(socketPath)
	if err != nil {
		return -1, fmt.Errorf("reading socket for CPU %s: %w", firstCPU, err)
	}

	socket, err := strconv.Atoi(strings.TrimSpace(string(socketData)))
	if err != nil {
		return -1, fmt.Errorf("parsing socket ID: %w", err)
	}

	return socket, nil
}
