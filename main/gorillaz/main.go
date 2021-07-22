package main

import (
	"apod/nasa"
	gorrilaz "github.com/skysoft-atm/gorillaz"
	http "net/http"
	"os"
)

func main() {

	context := &nasa.NasaContext{ApiKey: os.Args[1]}

	server := gorrilaz.New()
	server.Router.PathPrefix("/apod").Handler(http.HandlerFunc(context.GetData))

	<-server.Run()
	server.SetReady(true)

	select {}

}
