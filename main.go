package main

import (
	"log"
	"net"
	"net/http"

	"github.com/joshcarp/gop/gop/server"
)

func main() {
	http.HandleFunc("/", server.ServeHTTP)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(lis, nil))
}
