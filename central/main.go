package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
)

func configureInterfaces(interfaceName string) error {
	// Example tc command to set bandwidth limit to 1Mbps for all interfaces
	rand.Seed(time.Now().UnixNano())
	randomFloat := rand.Float64()
	s := strconv.FormatFloat(randomFloat*100, 'f', -1, 64)
	cmd := exec.Command("sudo", "tc", "qdisc", "change", "dev", interfaceName, "root", "netem", "loss", s)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error configuring interfaces: %v", err)
	}

	return nil
}

func main() {

	// Generate interface names
	for i := 0; i < 400; i++ {
		err := configureInterfaces(fmt.Sprintf("veth%d", i))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

}
