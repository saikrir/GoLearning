package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

func makeWebsiteCall(ctx context.Context, webSiteUrl string, resChan chan<- string, errChan chan<- error) {

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, webSiteUrl, nil)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Failed to make successful call to ", webSiteUrl)
		errChan <- err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		errChan <- err
	}

	resChan <- string(resBody)
}

const WAIT_DURATION = 4 * time.Second

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), WAIT_DURATION)
	defer cancel()

	websites := []string{
		"http://localhost:3000/",
		"http://www.amazon.com/",
	}

	resChan := make(chan string, len(websites))
	errChan := make(chan error, len(websites))

	for _, ws := range websites {
		go makeWebsiteCall(ctx, ws, resChan, errChan)
	}

	for i := 0; i < len(websites); i++ {
		select {
		case t := <-resChan:
			fmt.Println("Got Response ", t)
			cancel()
		case e := <-errChan:
			fmt.Println("Go an Error ", e)
			cancel()
		case <-ctx.Done():
			fmt.Println("Ok, times up, will cancel everything else")
			cancel()
		}
	}

	fmt.Println("Leaving with ", runtime.NumGoroutine())
}
