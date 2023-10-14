package chat

import (
	"GoChat/api"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"golang.org/x/net/websocket"
)

var unknownUsers = 0

type Server struct {
	conns map[*websocket.Conn]string

	messages []string
}

func newServer() *Server {
	return &Server{
		conns:    make(map[*websocket.Conn]string),
		messages: make([]string, 0),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	session, _ := api.Store.Get(ws.Request(), "session")
	login, ok := session.Values["login"]
	if !ok {
		login = "unknown_user_" + fmt.Sprint(unknownUsers)
		unknownUsers++
	}
	s.conns[ws] = login.(string)

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
		content := s.conns[ws] + ": " + parseMessage(msg)
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
	http.HandleFunc("/chatroom", getChatroomFunc(server))
}

func getChatroomFunc(server *Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !api.CheckSession(w, r) {
			return
		}
		tmpl := template.Must(template.ParseFiles("web/templates/chatroom.html"))
		tmpl.Execute(w, server.messages)
	}
}
