package chat

import (
	"golang.org/x/net/websocket"
)

type client struct {
	socket *websocket.Conn

	writeChan chan *string

	readChan chan *string

	room *room

	closed bool

	username string
}

func newClient(ws *websocket.Conn, room *room) *client {
	return &client{
		socket:    ws,
		readChan:  make(chan *string),
		writeChan: make(chan *string, 256),
		room:      room,
		closed:    false,
	}
}

func (c *client) read() {
	var msg string
	err := websocket.Message.Receive(c.socket, &msg)
	if err != nil {
		c.Close()
		return
	}
	c.readChan <- &msg
}

func (c *client) Write(msg string) {
	c.writeChan <- &msg
}

func (c *client) Close() {
	c.closed = true
	close(c.readChan)
	close(c.writeChan)
}

func (c *client) Listen() {
	go c.listenRead()
	go c.listenWrite()
}

func (c *client) listenRead() {
	for {
		if c.closed {
			return
		}
		go c.read()
		msg := <-c.readChan
		c.room.sendAll(msg)
	}

}

func (c *client) listenWrite() {
	for msg := range c.writeChan {
		websocket.Message.Send(c.socket, msg)
	}
}
