package service

import (
	"encoding/json"
	"example/kanban/appstate"
	db "example/kanban/db/sqlc"
	"log"
)

type server struct {
	id              uint32
	clients         map[uint32]*client
	clientCloseChan chan uint32
	clientNewChan   chan *client
	broadCastChan   chan uint8
}

func newServer() *server {
	return &server{
		id:              0,
		clients:         make(map[uint32]*client),
		clientCloseChan: make(chan uint32, 10),
		clientNewChan:   make(chan *client, 10),
		broadCastChan:   make(chan uint8),
	}
}

func (s *server) run() {
	for {
		select {
		case id := <-wsserver.clientCloseChan:
			delete(wsserver.clients, id)
		case client := <-wsserver.clientNewChan:
			wsserver.id += 1
			wsserver.clients[wsserver.id] = client
			go handleClient(client)
		case messageType := <-wsserver.broadCastChan:
			s.handleBroadCast(messageType)
		}
	}
}

func (s *server) handleBroadCast(messageType uint8) {

	switch messageType {
	case CreateTaskMessage:
		s.broadCastListTask()
	case DeleteTaskMessage:
		s.broadCastListTask()
	case UpdateTaskStatusMessage:
		s.broadCastListTask()
	default:
		log.Printf("service::handleBroadCast(): server receive unknown tag\n")
	}
}

func (s *server) broadCastListTask() {
	appState := appstate.GetAppState()
	ctx, cancelFunc := getCtx()
	defer cancelFunc()
	listTasks, err := db.ListTasks(ctx, appState.Db)
	if err != nil {
		log.Println("service::handleBroadCast(): failed to get list tasks: ", err)
		return
	}
	listTasksByte, err := json.Marshal(listTasks)
	if err != nil {
		log.Println("service::handleBroadCast(): failed to encode json: ", err)
		return
	}
	message := createMessageResponse(ListTaskUpdate, listTasksByte)
	for _, c := range s.clients {
		go c.sendMessage(message)
	}
}

func broadCastUpdate(updateType uint8) {
	go func() {
		wsserver.broadCastChan <- updateType
	}()
}
