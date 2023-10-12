package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool

	messages []string
}

func newServer() *Server {
	return &Server{
		conns:    make(map[*websocket.Conn]bool),
		messages: make([]string, 0),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 4192)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]
		content := parseMessage(msg)
		s.messages = append(s.messages, content)
		toSend := formatMessages(s)
		go s.broadcast(toSend)
	}
}

func parseMessage(msg []byte) string {
	decoded := make(map[string]any)
	err := json.Unmarshal(msg, &decoded)
	if err != nil {
		fmt.Println("unmarshalling error: ", err)
	}
	dec := decoded["chat_message"]
	return asString(dec)
}

func asString(dec any) string {
	outs, ok := dec.(string)
	if ok {
		return outs
	}
	outb, ok := dec.([]byte)
	if ok {
		return string(outb)
	}
	return ""
}

func formatMessages(s *Server) []byte {
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("web/templates/chatresponse.html"))
	if err := tmpl.Execute(&buf, s.messages); err != nil {
		return []byte("")
	}
	return buf.Bytes()
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}

		}(ws)
	}
}

func InitChat() {
	server := newServer()
	http.Handle("/chat", websocket.Handler(server.handleWS))
}
