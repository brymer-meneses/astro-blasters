package main

import (
	"flag"
	"log"

	"net"
	"net/http"
)

func main() {

	port := flag.String("port", "5500", "the port to host the server")

	flag.Parse()

	http.HandleFunc("/", http.FileServer(http.Dir("./static")).ServeHTTP)
	l, err := net.Listen("tcp", ":"+*port)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening at", l.Addr().(*net.TCPAddr).Port)
	http.Serve(l, nil)

}
