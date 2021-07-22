package main

import (
	"apod/nasa"
	"fmt"
	"net/http"
	"os"
)

func main() {

	apiKey := os.Args[1]
	//context := &nasa.Client{ApiKey: os.Args[1]}

	fmt.Printf("context start with APIKEY %s\n", apiKey)
	http.HandleFunc("/apod", nasa.Handler(apiKey))
	if err := http.ListenAndServe(":8081", nil); err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Printf("server stop")
}
