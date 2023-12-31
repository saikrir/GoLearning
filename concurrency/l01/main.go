package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func f01() {
	data := []int{1, 2, 3, 4}

	loopData := func(handleData chan<- int) {

		defer close(handleData)

		for _, i := range data {
			fmt.Println("Will write ", i)
			handleData <- i
		}
	}

	handleData := make(chan int, 5)
	go loopData(handleData)

	for val := range handleData {
		fmt.Println("Got ", val)
	}
}

func f02() {

	intStreamGen := func(nums ...int) <-chan int {

		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, i := range nums {
				fmt.Println("Sent ", i)
				intStream <- i
			}
		}()
		return intStream
	}

	intStream := intStreamGen(1, 2, 3, 4, 5)

	for i := range intStream {
		fmt.Println("Receieved ", i)
	}
}

func f03() {

	intStreamGen := func(nums ...int) <-chan int {

		intStream := make(chan int, len(nums))

		go func() {
			defer close(intStream)
			for _, i := range nums {
				fmt.Println("Sent ", i)
				intStream <- i
			}
		}()
		return intStream
	}

	consumer := func(intGen <-chan int) {
		for i := range intGen {
			fmt.Println("Inside Go Routine Received ", i)
		}
	}

	nums := []int{1, 2, 3, 4, 5}
	intStream := intStreamGen(nums...)
	consumer(intStream)
}

func f04() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		rBuffer := new(bytes.Buffer)

		for _, rbyte := range data {
			fmt.Fprintf(rBuffer, "%c", rbyte)
		}
		rBuffer.WriteTo(os.Stdout)
	}

	var pWg sync.WaitGroup

	pWg.Add(2)

	myData := []byte("Sai Katterishetty")

	go printData(&pWg, myData[3:])
	go printData(&pWg, myData[:3])
	pWg.Wait()
	fmt.Println("\nDone")
}

func f05() {

	generateRandoms := func(doneCh <-chan any) <-chan int {
		randStream := make(chan int)
		go func() {
			defer close(randStream)
			for {
				select {
				case <-doneCh:
					fmt.Println("Go Done Signal, Bye!")
					return
				case randStream <- rand.Intn(10):
				}
			}
		}()
		return randStream
	}

	doneStream := make(chan any)
	rStream := generateRandoms(doneStream)
	for i := 0; i < 5; i++ {
		fmt.Println("Got a number ", <-rStream)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	time.AfterFunc(1*time.Second, func() {
		defer wg.Done()
		close(doneStream)
	})
	wg.Wait()
	fmt.Println("Bye.. Exiting")
}

func main() {
	fmt.Println("Lesson L01")
	//f01()
	//f02()
	//f03()
	//f04()
	f05()
}
