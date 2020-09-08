package main

import (
	"fmt"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gop/retriever/retriever_proxy"
)

func main() {
	client := retriever_proxy.New(app.AppConfig{Proxy: "http://localhost:8080"})
	fmt.Println(client.Retrieve(`github.com/anz-bank/sysl-catalog`, `demo/demo.sysl`, `e6436737be76d167cd81ef69febecb6086a015bb`))
}
