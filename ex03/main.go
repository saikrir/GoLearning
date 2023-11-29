package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dollar float32

func (d dollar) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type webstore map[string]dollar

var myWebstore webstore

func (ws webstore) list(w http.ResponseWriter, req *http.Request) {
	for key, value := range ws {
		fmt.Fprintf(w, "%s costs %s \n", key, value)
	}
}

func (ws webstore) add(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	strItemCost := req.URL.Query().Get("value")
	cost, err := strconv.ParseFloat(strItemCost, 64)

	if err != nil {
		msg := fmt.Sprintf("Invalid input was provided, %s is not a valid prices \n", strItemCost)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if _, ok := ws[item]; ok {
		msg := fmt.Sprintf("Item %s already exists and will not be accepted \n", item)
		http.Error(w, msg, http.StatusConflict)
		return
	}
	ws[item] = dollar(cost)
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s-%f was added to db \n", item, cost)
}

func (ws webstore) update(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	strItemCost := req.URL.Query().Get("value")
	cost, err := strconv.ParseFloat(strItemCost, 64)

	if err != nil {
		msg := fmt.Sprintf("Invalid input was provided, %s is not a valid prices \n", strItemCost)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if _, ok := ws[item]; !ok {
		msg := fmt.Sprintf("Item %s does not exist, request will not be processed \n", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	ws[item] = dollar(cost)
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s-%f was Updated in db \n", item, cost)
}

func (ws webstore) delete(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")

	if _, ok := ws[item]; !ok {
		msg := fmt.Sprintf("Item %s does not exist, request will not be processed \n", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	delete(ws, item)
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s was Deleted in db \n", item)
}

func init() {
	myWebstore = map[string]dollar{
		"shoes": 50.00,
		"socks": 5.00,
	}
	fmt.Println("WebStore was initialized with item count ", len(myWebstore))
}

func main() {
	http.HandleFunc("/list", myWebstore.list)
	http.HandleFunc("/add", myWebstore.add)
	http.HandleFunc("/update", myWebstore.update)
	http.HandleFunc("/delete", myWebstore.delete)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
