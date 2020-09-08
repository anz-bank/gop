package main

import (
	"log"
	"net"
	"net/http"

	gop "github.com/joshcarp/gop"
)

func main() {
	http.HandleFunc("/", gop.ServeHTTP)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(lis, nil))
}
