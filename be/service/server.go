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
	data            *dataModel
}

type dataModel struct {
	listTask   []db.ListTasksRow
	listStatus []string
}

func newServer() *server {
	data := &dataModel{
		listTask:   nil,
		listStatus: nil,
	}
	s := &server{
		id:              0,
		clients:         make(map[uint32]*client),
		clientCloseChan: make(chan uint32, 10),
		clientNewChan:   make(chan *client, 10),
		broadCastChan:   make(chan uint8),
		data:            data,
	}
	if _, err := s.getListTask(true); err != nil {
		log.Fatal("server::newServer(): failed to update list task")
	}
	if _, err := s.getListStatus(true); err != nil {
		log.Fatal("server::newServer(): failed to update list status")
	}
	return s
}

func (s *server) run() {
	for {
		select {
		case id := <-wsserver.clientCloseChan:
			delete(wsserver.clients, id)
		case client := <-wsserver.clientNewChan:
			wsserver.id += 1
			wsserver.clients[wsserver.id] = client
			go func() {
				if msg, err := s.getListStatus(false); err == nil {
					client.sendMessage(msg)
				} else {
					log.Println("server::run(): failed to get list status")
				}
				if msg, err := s.getListTask(false); err == nil {
					client.sendMessage(msg)
				} else {
					log.Println("server::run(): failed to get list status")
				}
				handleClient(client)
			}()
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
	message, err := s.getListTask(true)
	if err != nil {
		return
	}
	for _, c := range s.clients {
		go c.sendMessage(message)
	}
}

func (s *server) getListTask(needUpdate bool) ([]byte, error) {
	if needUpdate {
		appState := appstate.GetAppState()
		ctx, cancelFunc := getCtx()
		defer cancelFunc()
		listTasks, err := db.ListTasks(ctx, appState.Db)
		if err != nil {
			log.Println("service::getListTaskMessage(): failed to get list tasks: ", err)
			return nil, err
		}
		s.data.listTask = listTasks
	}
	listTasksByte, err := json.Marshal(s.data.listTask)
	if err != nil {
		log.Println("service::getListTaskMessage(): failed to encode json: ", err)
		return nil, err
	}
	message := createMessageResponse(ListTaskUpdate, listTasksByte)
	return message, nil
}

func (s *server) getListStatus(needUpdate bool) ([]byte, error) {
	if needUpdate {
		appState := appstate.GetAppState()
		ctx, cancelFunc := getCtx()
		defer cancelFunc()
		listStatus, err := db.ListStatus(ctx, appState.Db)
		if err != nil {
			log.Println("service::getListStatus(): failed to get list tasks: ", err)
			return nil, err
		}
		s.data.listStatus = listStatus
	}
	listStatusByte, err := json.Marshal(s.data.listStatus)
	if err != nil {
		log.Println("service::getListStatus(): failed to encode json: ", err)
		return nil, err
	}
	message := createMessageResponse(ListStatusUpdate, listStatusByte)
	return message, nil
}

func broadCastUpdate(updateType uint8) {
	go func() {
		wsserver.broadCastChan <- updateType
	}()
}
