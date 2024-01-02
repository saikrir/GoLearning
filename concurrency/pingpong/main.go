package main

import (
	"fmt"
	"time"
)

func pong(pingCh <-chan int, pongCh chan<- int, doneCh <-chan any) {

	var lastVal int

	for {
		select {
		case <-doneCh:
			fmt.Println("Done")
			return
		case lastVal = <-pingCh:
			fmt.Println("Ping ", lastVal)
			time.Sleep(1 * time.Second)
			pongCh <- lastVal + 1
		}
	}
}

func main() {
	fmt.Println("PingPong")
	var initVal int
	pingCh := make(chan int, 1)
	pongCh := make(chan int, 1)
	doneCh := make(chan any)
	defer close(doneCh)
	defer close(pingCh)
	defer close(pongCh)
	go func() {
		pingCh <- initVal
	}()

	time.AfterFunc(5*time.Second, func() {
		doneCh <- "Done"
	})

	go pong(pingCh, pongCh, doneCh)

loop:
	for {
		select {
		case <-doneCh:
			break loop

		case initVal = <-pongCh:
			fmt.Println("Pong ", initVal)
			time.Sleep(1 * time.Second)
			pingCh <- initVal + 1
		}
	}

}
