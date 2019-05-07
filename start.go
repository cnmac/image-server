package main

import (
	"github.com/cnmac/image-server/httphandler"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.POST("/upload", httphandler.Upload)
	router.GET("/download/:name", httphandler.Download)
	log.Fatal(http.ListenAndServe(":8080", router))
}
