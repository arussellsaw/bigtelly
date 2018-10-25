package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var (
	update chan urlMessage

	mu       sync.Mutex
	current  string
	urls     []string
	interval time.Duration
)

func main() {
	update = make(chan urlMessage, 1)
	update <- urlMessage{URLs: []string{"https://yaas.cat"}, Interval: 5 * time.Minute}
	go trackURLS()
	r := mux.NewRouter()
	r.Handle("/url", http.HandlerFunc(urlHandler))
	r.Handle("/currenturl", http.HandlerFunc(currentURLHandler))
	ctx := context.Background()
	err := initChrome(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080", r)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	u := urlMessage{}
	err := dec.Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	update <- u
}

type urlMessage struct {
	URLs     []string
	Interval time.Duration
}

func currentURLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CURRENT URL", currentURL())
	w.Write([]byte(currentURL()))
}

func currentURL() string {
	mu.Lock()
	defer mu.Unlock()
	return current
}

func trackURLS() {
	var (
		interval time.Duration = 10 * time.Second
		urls     []string
		i        int
	)
L:
	for {
		t := time.After(interval)
		select {
		case um := <-update:
			if len(um.URLs) <= 0 {
				fmt.Println("cannot set empty urls")
				continue L
			}
			interval = um.Interval
			urls = um.URLs
			current = urls[0]
			fmt.Printf("updated urls to %s, interval to %s\n", urls, interval)
			continue L
		case <-t:
			i++
			if i == len(urls) {
				i = 0
			}
			mu.Lock()
			fmt.Println("set url to", urls[i])
			current = urls[i]
			mu.Unlock()
		}
	}
}
