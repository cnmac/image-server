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
	router.GET("/download/:alias", httphandler.Download)
	router.GET("/delete/:alias", httphandler.Delete)
	log.Fatal(http.ListenAndServe(":9999", router))
}
