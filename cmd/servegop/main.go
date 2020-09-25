package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

//func main() {
//	http.HandleFunc("/", gop.ServeHTTP)
//	lis, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Fatal(http.Serve(lis, nil))
//}

func main() {
	a := map[string]string{"github.com/anz-bank/sysl@1234": "1", "github.com/anz-bank/sysl/": "2"}
	b, _ := yaml.Marshal(a)
	fmt.Println(string(b))

}
