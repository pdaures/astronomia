package main

import (
	"apod/nasa"
	"fmt"
	http "net/http"
	"os"
)

func main() {

	context := &nasa.NasaContext{ApiKey: os.Args[1]}

	fmt.Printf("context start with APIKEY %s\n", context.ApiKey)
	http.HandleFunc("/apod", context.GetData)
	http.ListenAndServe(":8081", nil)
	fmt.Printf("server stop")

}
