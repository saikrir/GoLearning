package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const baseUrl = "https://xkcd.com/%d/info.0.json"
const fileName = "xkcd.json"

type xkcdItem struct {
	Year       string `json:year`
	Month      string `json:month`
	Day        string `json:day`
	Title      string `json:title`
	Transcript string `json:transcript`
	ImageUrl   string `json:img`
}

func downloadApiResponse(aNumber int) (resp []byte, err error) {

	currentUrl := fmt.Sprintf(baseUrl, aNumber)
	response, err := http.Get(currentUrl)

	if err != nil {
		fmt.Println("Unable to fetch data ", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("An unacceptable status code was returned [%d]", response.StatusCode)
		return nil, errors.New(errorMessage)
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func init() {
	fmt.Println("Starting Program")
}

func loadData() {
	var (
		errorCnt     int
		completedCnt int
		rawResponse  []byte
		err          error
		output       io.WriteCloser
	)

	if output, err = os.Create(fileName); err != nil {
		log.Fatalln("Failed to open file ", fileName)
	}

	fmt.Fprint(output, "[")
	defer fmt.Fprint(output, "]")

	for i := 1; errorCnt < 2; i++ {

		if rawResponse, err = downloadApiResponse(i); err != nil {
			fmt.Println("Failed to load data ", err)
			errorCnt++
			continue
		}
		if completedCnt > 0 {
			fmt.Fprint(output, ",")
		}

		if _, err := io.Copy(output, bytes.NewBuffer(rawResponse)); err != nil {
			log.Fatalln("Failed to write data ", err)
		}
		fmt.Printf("Completed downloading %d \n", completedCnt)
		completedCnt++
	}
}

func searchData(searchTerms string) (results []xkcdItem) {
	var (
		err      error
		items    []xkcdItem
		retItems []xkcdItem
		xkcdFile *os.File
	)

	if xkcdFile, err = os.Open(fileName); err != nil {
		log.Fatalln("Failed to load file", err)
	}
	defer xkcdFile.Close()

	if err = json.NewDecoder(xkcdFile).Decode(&items); err != nil {
		log.Fatalln("Failed to Parse JSON")
	}

	fmt.Printf("Loaded %d items \n", len(items))

	terms := strings.Fields(strings.ToLower(searchTerms))

outer:
	for _, xkcdItem := range items {
		lowerTitle := strings.ToLower(xkcdItem.Title)
		lowerTranscript := strings.ToLower(xkcdItem.Transcript)
		for _, term := range terms {
			if !strings.Contains(lowerTitle, term) && !strings.Contains(lowerTranscript, term) {
				continue outer
			}
			retItems = append(retItems, xkcdItem)
		}
	}

	return retItems
}

func main() {
	//loadData()
	fmt.Printf("Found %d \n", len(searchData("someone bed sleep")))
}
