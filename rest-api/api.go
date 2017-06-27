package rest_api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"log"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Sse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	broker := NewServer()
	broker.ServeHTTP(w, r)
	go func() {
		for {
			time.Sleep(time.Second)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()
}
