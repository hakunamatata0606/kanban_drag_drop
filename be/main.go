package main

import (
	"example/kanban/appstate"
	"example/kanban/service"
	"net/http"
)

func main() {
	appState := appstate.GetAppState()
	http.HandleFunc("/ws", service.ServeWs)
	http.ListenAndServe(appState.Config.ServerUrl, nil)
}
