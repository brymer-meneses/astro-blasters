package main

import (
	"flag"
	"log"
	"space-shooter/server"
)

func main() {
	port := flag.String("port", "5500", "the port to host the server")
	flag.Parse()

	server := server.NewServer()

	if err := server.Start(*port); err != nil {
		log.Fatal(err)
	}

}
