package main

import (
	"github.com/xyalan/migateway-api"
)

func main() {
	_, err := migateway.NewMiHomeManager(nil)
	if err != nil {
		panic(err)
	}

	//do something...
	//time.Sleep(10 * time.Second)
	select {}
}
