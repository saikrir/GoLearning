package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Consumer will wait for producer to feed it input
func waitForTask() {
	myChannel := make(chan string)

	go func() {
		val := <-myChannel
		fmt.Println("Op:", val)
	}()
	delay := rand.Intn(5) + 1
	time.Sleep(time.Duration(delay) * time.Second)
	myChannel <- fmt.Sprintf("Task Fed after waiting %d seconds", delay)
	time.Sleep(time.Second)
	defer close(myChannel)
}

// Consumer already has input and will let Producer know once done with Result
func waitForResult() {
	myChannel := make(chan string)

	go func() {
		delay := rand.Intn(5) + 1
		time.Sleep(time.Duration(delay) * time.Second)
		myChannel <- fmt.Sprintf("Produced after waiting %d", delay)
	}()

	val := <-myChannel
	fmt.Println("Result: ", val)
	defer close(myChannel)
}

// Conusmer already has input and will let producer know once done by closing channel
func waitForFinished() {
	myChannel := make(chan struct{})
	go func() {
		delay := rand.Intn(5) + 1
		time.Sleep(time.Duration(delay) * time.Second)
		fmt.Println("Finished my work ", delay)
		close(myChannel)
	}()

	_, result := <-myChannel

	if !result {
		fmt.Println("Task Completed")
	}
}

// There are pool of consumers that are waiting for input via queue from Produer
func pooling() {
	workQueue := make(chan string)
	startTime := time.Now()

	worker := func(num int) {
		for work := range workQueue {
			time.Sleep(1 * time.Second)
			fmt.Printf("Completed work by [%d] : [%s] \n", num, strings.ToUpper(work))
		}
		fmt.Printf("Worker [%d] will shutdown \n", num)
	}

	poolSize := 4
	for i := 0; i < poolSize; i++ {
		go worker(i)
	}

	tasks := []string{"Jaya", "Vicky", "Shiva", "Dinesh", "Nitin", "Sai", "Jagan", "Narasimha", "Mangala", "Shanta", "Ishita", "Nymisha"}

	for _, task := range tasks {
		workQueue <- task
	}

	close(workQueue)
	completionDuration := time.Since(startTime).Round(time.Second)
	fmt.Printf("%d tasks Completed in %d seconds \n", len(tasks), completionDuration)
}

// multiple Consumers have input and will use buffered channel, producer will count down for all tasks to complete

func fanOut() {

	numWorkers := 10

	workQueue := make(chan string, numWorkers)

	worker := func(workerN int) {
		delay := time.Duration(rand.Intn(5) + 1)
		time.Sleep(delay * time.Second)
		workQueue <- fmt.Sprintf("Worker %d, completed after %d seconds ", workerN, delay)
	}

	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}

	fmt.Println("All Workers Launched")

	totalExpectedResults := numWorkers

	for totalExpectedResults > 0 {
		result := <-workQueue
		fmt.Printf("[Result] : %s, [%d] pending \n", result, totalExpectedResults)
		totalExpectedResults--
	}

	time.Sleep(time.Second)
	fmt.Println("Completed Fanout")
}

// consumer will pull in work, producer will try to feed input thru buffered channel, if it blocks on feed, then it will drop it.
func drop() {
	workQueue := make(chan string, 2)

	for i := 0; i < 2; i++ {
		go func() {
			for workItem := range workQueue {
				delay := time.Duration(rand.Intn(5) + 1)
				time.Sleep(delay * time.Second)
				fmt.Println("Completed ", workItem)
			}
		}()
	}

	for i := 0; i < 10; i++ {
		taskName := fmt.Sprintf("Task %d", i)
		select {
		case workQueue <- taskName:
		default:
			fmt.Println("Will Drop ", taskName)
		}
		time.Sleep(time.Second)
	}

	fmt.Println("All Done")
}

// consumer will perform a task, producer waits for sometime and cancels if task does not complete on Deadline
func cancellation() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	myQueue := make(chan string, 1)

	go func() {
		delay := time.Duration(rand.Intn(5) + 1)
		fmt.Println("Delay ", delay)
		time.Sleep(delay * time.Second)
		myQueue <- "Completed"
	}()

	select {
	case v := <-myQueue:
		fmt.Println("Value ", v)
	case <-ctx.Done():
		fmt.Println("Will Cancel ")
	}
	fmt.Println("All done")
}

func main() {
	//waitForTask()
	//waitForResult()
	//waitForFinished()
	//pooling()
	//fanOut()
	//drop()
	cancellation()
}
