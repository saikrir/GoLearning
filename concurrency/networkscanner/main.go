package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type pingResult struct {
	host      string
	reachable bool
}

func netScanner(netPrefix string, results chan<- pingResult) {
	for i := 1; i < 255; i++ {
		ipAddress := fmt.Sprintf("%s.%d", netPrefix, i)
		go func() {
			out, _ := exec.Command("ping", "-c 1", ipAddress).Output()
			//fmt.Println("Output ", ipAddress, string(out))
			hostReachable := strings.Contains(string(out), "1 packets transmitted, 1 packets received")
			results <- pingResult{host: ipAddress, reachable: hostReachable}
		}()
	}
}

func main() {
	fmt.Println("Starting Network Scanner")
	results := make(chan pingResult, 255)
	go netScanner("192.168.86", results)
	for i := 1; i < 255; i++ {
		result := <-results
		if result.reachable {
			fmt.Println("Found ", result.host)
		}
	}
	close(results)
	fmt.Println("Bye ", runtime.NumGoroutine())
}
