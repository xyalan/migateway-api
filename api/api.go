package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Sse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	broker := NewServer()
	flusher, messageChan := broker.ServeHTTP(w, r)
	go func() {
		for {
			time.Sleep(time.Second)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()

	for {
		log.Println("fffff")
		fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
		flusher.Flush()
	}
}
