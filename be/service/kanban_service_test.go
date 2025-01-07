package service_test

import (
	"encoding/json"
	db "example/kanban/db/sqlc"
	"example/kanban/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func prepareCreateMessage(t *testing.T, name string, title string, descr string, status string) []byte {
	request := db.CreateTaskRequest{
		Name:        name,
		Title:       title,
		Description: descr,
		Status:      status,
	}
	bytes, err := json.Marshal(&request)
	assert.Nil(t, err)
	var msg []byte
	msg = append(msg, service.CreateTaskMessage)
	msg = append(msg, bytes...)
	return msg
}

func prepareDeleteMessage(t *testing.T, name string) []byte {
	request := db.DeleteTaskRequest{
		Name: name,
	}
	bytes, err := json.Marshal(&request)
	assert.Nil(t, err)
	var msg []byte
	msg = append(msg, service.DeleteTaskMessage)
	msg = append(msg, bytes...)
	return msg
}

func prepareUpdateTaskStatusMessage(t *testing.T, name string, status string) []byte {
	request := db.UpdateTaskRequest{
		Name:   name,
		Status: status,
	}
	bytes, err := json.Marshal(&request)
	assert.Nil(t, err)
	var msg []byte
	msg = append(msg, service.UpdateTaskStatusMessage)
	msg = append(msg, bytes...)
	return msg
}

func TestKanbanServiceSimple(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(service.ServeWs))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)
	msg := prepareCreateMessage(t, "test1", "title1", "desc1", "idea")
	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	assert.Nil(t, err)

	msgtype, msg, err := ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAck, int(msg[0]))
	msg = msg[1:]
	var resp service.MessageAckResponse
	err = json.Unmarshal(msg, &resp)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAckResponse{Status: http.StatusOK, Message: ""}, resp)

	msgtype, msg, err = ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.ListTaskUpdate, int(msg[0]))
	msg = msg[1:]
	var resp1 []db.ListTasksRow
	err = json.Unmarshal(msg, &resp1)
	assert.Nil(t, err)
	assert.Equal(t, []db.ListTasksRow{db.ListTasksRow{Name: "test1", Title: "title1", Description: "desc1", Status: "idea"}}, resp1)

	msg = prepareUpdateTaskStatusMessage(t, "test1", "done")
	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	assert.Nil(t, err)

	msgtype, msg, err = ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAck, int(msg[0]))
	msg = msg[1:]
	err = json.Unmarshal(msg, &resp)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAckResponse{Status: http.StatusOK, Message: ""}, resp)

	msgtype, msg, err = ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.ListTaskUpdate, int(msg[0]))
	msg = msg[1:]
	err = json.Unmarshal(msg, &resp1)
	assert.Nil(t, err)
	assert.Equal(t, []db.ListTasksRow{db.ListTasksRow{Name: "test1", Title: "title1", Description: "desc1", Status: "done"}}, resp1)

	msg = prepareDeleteMessage(t, "test1")
	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	assert.Nil(t, err)

	msgtype, msg, err = ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAck, int(msg[0]))
	msg = msg[1:]
	err = json.Unmarshal(msg, &resp)
	assert.Nil(t, err)
	assert.Equal(t, service.MessageAckResponse{Status: http.StatusOK, Message: ""}, resp)

	msgtype, msg, err = ws.ReadMessage()
	assert.Equal(t, msgtype, websocket.BinaryMessage)
	assert.Nil(t, err)
	assert.Equal(t, service.ListTaskUpdate, int(msg[0]))
	msg = msg[1:]
	err = json.Unmarshal(msg, &resp1)
	assert.Nil(t, err)
	assert.Equal(t, []db.ListTasksRow(nil), resp1)
}
