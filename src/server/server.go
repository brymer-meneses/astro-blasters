package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Server struct {
	serveMux http.ServeMux
}

func NewServer() *Server {
	cs := &Server{}
	cs.serveMux.HandleFunc("/events/ws", cs.ws)
	cs.serveMux.HandleFunc("/", cs.root)
	return cs
}

func (self *Server) Start(port string) error {
	log.Printf("Listening at %s", port)
	return http.ListenAndServe(":"+port, &self.serveMux)
}

func (self *Server) root(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func (self *Server) ws(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "Connection Failed")
	}
	defer c.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var v interface{}
	if err := wsjson.Read(ctx, c, &v); err != nil {
		log.Fatal(err)
	}

	c.Close(websocket.StatusNormalClosure, "")
}
