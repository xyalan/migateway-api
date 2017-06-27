package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"github.com/xyalan/migateway-api/rest-api"
)

func main() {
	router := httprouter.New()
	router.GET("/", rest_api.Index)
	router.GET("/sse", rest_api.Sse)
	router.GET("/hello/:name", rest_api.Hello)
	log.Fatal(http.ListenAndServe(":8080", router))
}