package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

func configureInterfaces(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Loop until the channel is closed
	for interfaceName := range ch {
		// Example tc command to set bandwidth limit to 1Mbps for each interface
		cmd := exec.Command("sudo", "tc", "qdisc", "change", "dev", interfaceName, "root", "netem", "loss", "1%")

		// Execute the command
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error configuring interface %s: %v\n", interfaceName, err)
			continue
		}

		//fmt.Printf("Interface %s configured successfully\n", interfaceName)
	}
}

func main() {
	// List of interface names
	interfaceNames := make(chan string, 10000)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numGoroutines := 10

	// Add 20 tasks to WaitGroup
	wg.Add(numGoroutines)

	// Start 20 goroutines to configure the interfaces concurrently
	for i := 0; i < numGoroutines; i++ {
		go configureInterfaces(interfaceNames, &wg)
	}

	// Generate interface names and send them to the channel
	for i := 0; i < 400; i++ {
		rand.Seed(time.Now().UnixNano())
		// 生成0到19之间的随机整数
		//randomNumber := rand.Intn(20)
		//interfaceNames <- fmt.Sprintf("veth%d", randomNumber)
		interfaceNames <- fmt.Sprintf("veth%d", i)
	}

	// Close the channel to signal that no more interface names will be sent
	close(interfaceNames)

	// Wait for all goroutines to finish
	wg.Wait()

	//fmt.Println("All interfaces configured successfully")
}
