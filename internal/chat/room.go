package chat

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type room struct {
	messages []*string

	clients map[*client]bool

	joinChan chan *client

	leaveChan chan *client

	forward chan *string
}

func newRoom() *room {
	return &room{
		clients:   make(map[*client]bool),
		joinChan:  make(chan *client),
		leaveChan: make(chan *client),
		forward:   make(chan *string),
	}
}

func (r *room) Join(c *client) {
	r.joinChan <- c
}

func (r *room) Leave(c *client) {
	r.leaveChan <- c
}

func (r *room) SendAll(msg *string) {
	r.forward <- msg
}

func (r *room) sendPast(c *client) {
	go func() {
		for _, msg := range r.messages {
			c.Write(*msg)
		}
	}()
}

func (r *room) sendAll(msg *string) {
	for c, _ := range r.clients {
		go c.Write(*msg)
	}
}

func (r *room) Handler() http.Handler {
	return websocket.Handler(func(c *websocket.Conn) {
		defer func() {
			if err := c.Close(); err != nil {
				log.Println("Error: ", err)
			}
		}()
		client := newClient(c, r)
		r.Join(client)
		client.Listen()
	})
}

func (r *room) run() {
	for {
		select {
		case c := <-r.joinChan:
			r.clients[c] = true
			msg := "New client joined!"
			r.sendAll(&msg)
			r.sendPast(c)
		case c := <-r.leaveChan:
			delete(r.clients, c)
			msg := "Client left :("
			r.sendAll(&msg)
		case msg := <-r.forward:
			r.messages = append(r.messages, msg)
			r.sendAll(msg)
		}
	}
}
