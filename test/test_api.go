package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/xyalan/migateway-api/restapi"
)

func main() {
	router := httprouter.New()
	router.GET("/", api.Index)
	router.GET("/sse", api.Sse)
	router.GET("/hello/:name", api.Hello)
	log.Fatal(http.ListenAndServe(":8080", router))
}