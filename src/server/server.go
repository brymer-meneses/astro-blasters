package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		_, sent, err := c.Read(ctx)
		status := websocket.CloseStatus(err)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			break
		}

		if status == websocket.StatusGoingAway || status == websocket.StatusAbnormalClosure {
			break
		}

		if err != nil {
			break
		}

		log.Printf("Message %s\n", string(sent))
	}

}
