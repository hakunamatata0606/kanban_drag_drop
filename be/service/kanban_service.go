package service

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsserver *server
var upgrader = &websocket.Upgrader{}

func init() {
	wsserver = newServer()
	go wsserver.run()
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to get ws conn: ", err)
		return
	}
	log.Printf("service::ServeWs(): new ws connection - %s\n", r.RemoteAddr)
	addClient(conn)
}
