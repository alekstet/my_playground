package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, playground")
	})

	fmt.Println("Starting server...")
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		fmt.Println(http.Serve(l, nil))
	}()

	fmt.Println("Sending request...")
	res, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Reading response...")
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		fmt.Println(err)
	}
}
