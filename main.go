package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
        "os"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), proxy))
}
