package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type pingResult struct {
	host      string
	reachable bool
}

func netScanner(netPrefix string, results chan<- pingResult) {
	for i := 1; i < 255; i++ {
		ipAddress := fmt.Sprintf("%s.%d", netPrefix, i)
		go func(ipAddress string) {
			out, _ := exec.Command("ping", "-c 1", ipAddress).Output()
			//fmt.Println("Output ", ipAddress, string(out))
			hostReachable := strings.Contains(string(out), fmt.Sprintf("bytes from %s", ipAddress))
			results <- pingResult{host: ipAddress, reachable: hostReachable}
		}(ipAddress)
	}
}

func main() {
	fmt.Println("Starting Network Scanner")
	startTime := time.Now()
	results := make(chan pingResult, 255)
	go netScanner("192.168.86", results)
	nFound := 0
	for i := 1; i < 255; i++ {
		result := <-results
		if result.reachable {
			nFound++
			fmt.Println("Found ", result.host)
		}
	}
	close(results)
	fmt.Println("Discovered ", nFound, time.Since(startTime).Round(time.Second))
}
