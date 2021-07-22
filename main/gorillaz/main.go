package main

import (
	"apod/nasa"
	"os"

	gorrilaz "github.com/skysoft-atm/gorillaz"
)

func main() {

	server := gorrilaz.New()

	server.Router.HandleFunc("/apod", nasa.CacheHandler(&nasa.Store{Path: "./"}, nasa.Handler(os.Args[1])))

	<-server.Run()
	server.SetReady(true)

	select {}

}
