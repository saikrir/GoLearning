package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)


const REQ_TO = 3

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got a request from ", r.Header.Get("User-Agent"))
	time.Sleep(REQ_TO * time.Second)
	fmt.Fprintf(w, "Hello World")
}

func main() {
	http.HandleFunc("/", sayHello)
	log.Fatalln(http.ListenAndServe(":3000", nil))
}
