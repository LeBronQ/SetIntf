package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Task struct {
	interfaceName string
	wg            *sync.WaitGroup
}

func (t *Task) configureInterfaces() {
	// Example tc command to set bandwidth limit to 1Mbps for each interface
	cmd := exec.Command("sudo", "tc", "qdisc", "change", "dev", t.interfaceName, "root", "netem", "loss", "1%")

	// Execute the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error configuring interface %s: %v\n", t.interfaceName, err)
	}
	t.wg.Done()
}

func main() {
	runtime := 1000
	defer ants.Release()
	var wg sync.WaitGroup
	taskFunc := func(data interface{}) {
		task := data.(*Task)
		task.configureInterfaces()
	}
	p, _ := ants.NewPoolWithFunc(10, taskFunc)
	defer p.Release()
	// Generate interface names and send them to the channel
	wg.Add(runtime)
	for i := 0; i < runtime; i++ {
		rand.Seed(time.Now().UnixNano())
		// 生成0到19之间的随机整数
		randomNumber := rand.Intn(20)
		name := fmt.Sprintf("veth%d", randomNumber)
		task := &Task{
			interfaceName: name,
			wg:            &wg,
		}
		p.Invoke(task)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
}
