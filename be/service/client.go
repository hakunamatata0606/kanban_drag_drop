package service

import (
	"context"
	"encoding/json"
	"example/kanban/appstate"
	db "example/kanban/db/sqlc"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	CreateTaskMessage = iota
	DeleteTaskMessage
	UpdateTaskStatusMessage
	GetListTask
	MessageAck
	ListTaskUpdate
)

type client struct {
	id      uint32
	stopped bool
	conn    *websocket.Conn
}

type MessageAckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func createMessageResponse(messageType uint8, payload []byte) []byte {
	var message []byte
	message = append(message, messageType)
	message = append(message, payload...)
	return message
}

func createMessageAck(status int, message string) []byte {
	msgJson := MessageAckResponse{
		Status:  status,
		Message: message,
	}
	bytes, err := json.Marshal(&msgJson)
	if err != nil {
		log.Fatal("createMessageAck(): Failed to encode json: ", err)
	}
	return createMessageResponse(MessageAck, bytes)
}

func handleClient(c *client) {
	for {
		if c.stopped {
			return
		}
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("service::handleClient(): client[%d] error read message: %s\n", c.id, err)
			c.close()
			return
		}
		if mt == websocket.CloseMessage {
			log.Printf("service::handleClient(): client[%d] closed connection\n", c.id)
			c.close()
			return
		}
		c.handleClientMessage(message)
	}
}

func (c *client) sendMessage(msg []byte) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Printf("service::sendMessage(): client[%d] failed to send message\n", c.id)
		c.close()
	}
}

func (c *client) handleClientMessage(msg []byte) {
	if len(msg) == 0 {
		log.Printf("service::handleClientMessage(): client[%d] receive empty message\n", c.id)
		return
	}
	tag := uint8(msg[0])
	payload := msg[1:]
	switch tag {
	case CreateTaskMessage:
		log.Printf("service::handleClientMessage(): client[%d] receive create task tag\n", c.id)
		c.handleCreateTaskRequest(payload)
	case DeleteTaskMessage:
		log.Printf("service::handleClientMessage(): client[%d] receive delete task tag\n", c.id)
		c.handleDeleteTaskRequest(payload)
	case UpdateTaskStatusMessage:
		log.Printf("service::handleClientMessage(): client[%d] receive update task status tag\n", c.id)
		c.handleUpdateTaskRequest(payload)
	default:
		log.Printf("service::handleClientMessage(): client[%d] receive unknown tag\n", c.id)
		ack := createMessageAck(http.StatusBadRequest, "")
		c.sendMessage(ack)
	}
}

func (c *client) handleCreateTaskRequest(msg []byte) {
	appState := appstate.GetAppState()

	var createRequest db.CreateTaskRequest
	err := json.Unmarshal(msg, &createRequest)
	if err != nil {
		log.Printf("service::handleCreateTaskReques(): client[%d] fail to unmarshal json\n", c.id)
		ack := createMessageAck(http.StatusBadRequest, "")
		c.sendMessage(ack)
		return
	}
	ctx, cancelFunc := getCtx()
	defer cancelFunc()
	err = db.CreateTask(ctx, appState.Db, &createRequest)
	if err != nil {
		log.Printf("service::handleCreateTaskReques(): client[%d] fail to create task - %s\n", c.id, err)
		ack := createMessageAck(http.StatusInternalServerError, "")
		c.sendMessage(ack)
		return
	}
	ack := createMessageAck(http.StatusOK, "")
	c.sendMessage(ack)
	broadCastUpdate(CreateTaskMessage)
}

func (c *client) handleUpdateTaskRequest(msg []byte) {
	appState := appstate.GetAppState()

	var updateTaskRequest db.UpdateTaskRequest
	err := json.Unmarshal(msg, &updateTaskRequest)
	if err != nil {
		log.Printf("service::handleUpdateTaskRequest(): client[%d] fail to unmarshal json\n", c.id)
		ack := createMessageAck(http.StatusBadRequest, "")
		c.sendMessage(ack)
		return
	}

	ctx, cancelFunc := getCtx()
	defer cancelFunc()
	err = db.UpdateTask(ctx, appState.Db, &updateTaskRequest)
	if err != nil {
		log.Printf("service::handleUpdateTaskRequest(): client[%d] fail to update task - %s\n", c.id, err)
		ack := createMessageAck(http.StatusInternalServerError, "")
		c.sendMessage(ack)
		return
	}
	ack := createMessageAck(http.StatusOK, "")
	c.sendMessage(ack)
	broadCastUpdate(UpdateTaskStatusMessage)
}

func (c *client) handleDeleteTaskRequest(msg []byte) {
	appState := appstate.GetAppState()

	var deleteRequest db.DeleteTaskRequest
	err := json.Unmarshal(msg, &deleteRequest)
	if err != nil {
		log.Printf("service::handleDeleteTaskRequest(): client[%d] fail to unmarshal json\n", c.id)
		ack := createMessageAck(http.StatusBadRequest, "")
		c.sendMessage(ack)
		return
	}
	ctx, cancelFunc := getCtx()
	defer cancelFunc()
	err = db.DeleteTask(ctx, appState.Db, &deleteRequest)
	if err != nil {
		log.Printf("service::handleDeleteTaskRequest(): client[%d] fail to delete task - %s\n", c.id, err)
		ack := createMessageAck(http.StatusInternalServerError, "")
		c.sendMessage(ack)
		return
	}
	ack := createMessageAck(http.StatusOK, "")
	c.sendMessage(ack)
	broadCastUpdate(DeleteTaskMessage)
}

func getCtx() (context.Context, context.CancelFunc) {
	appState := appstate.GetAppState()
	return context.WithTimeout(context.Background(), time.Duration(*appState.Config.QueryTimeout)*time.Second)
}

func (c *client) close() {
	c.stopped = true
	c.conn.Close()
	wsserver.clientCloseChan <- c.id
}

func addClient(conn *websocket.Conn) {
	client := &client{
		id:      0,
		stopped: false,
		conn:    conn,
	}
	wsserver.clientNewChan <- client
}
