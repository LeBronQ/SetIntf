package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

func configureInterfaces(wg *sync.WaitGroup) error {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	randomFloat := rand.Float64()
	s := strconv.FormatFloat(randomFloat*100, 'f', -1, 64)
	cmd := exec.Command("sudo", "tc", "qdisc", "change", "dev", "lo", "root", "netem", "loss", s)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error configuring interfaces: %v", err)
	}

	return nil
}

func main() {

	// Generate interface names
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go configureInterfaces(&wg)
	}
	wg.Wait()

}
