package webserver

import (
	"crm-backend/internal"
	"crm-backend/types"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type SocketSession struct {
	Context *types.WebSocketContext
}

type SocketManager struct {
	controllers map[string]map[string]func(context *types.RequestContext)
	Sessions    []*SocketSession
}

func (sm *SocketManager) log(logMessage string) {
	fmt.Println(logMessage)
}

func (sm *SocketManager) Updater(Controller string, Data []interface{}) {
	if len(Data) == 0 {
		return
	}

	for _, session := range sm.Sessions {
		if session.Context.LastController == Controller {
			packet := &types.WebSocketPackage{
				Controller: Controller,
				Action:     "Update",
				Serial:     GetFunnyWord(),
				Error:      "",
				Data:       Data,
			}
			sm.WriteMessage(session.Context, packet)
		}
	}
}

func (sm *SocketManager) WriteMessage(wsCtx *types.WebSocketContext, packet *types.WebSocketPackage) {
	jsonDate, err := internal.InsecureMarshal(packet.Data)
	if err != nil {
		sm.log(err.Error())
	}
	packet.Hash = internal.GetHash(jsonDate)

	jsonStructure, err := internal.InsecureMarshal(packet)
	if err != nil {
		sm.log(err.Error())
	}

	err = wsCtx.Socket.WriteMessage(websocket.TextMessage, jsonStructure)
	if err != nil {
		sm.log(err.Error())
	}
}

func (sm *SocketManager) do(context *types.WebSocketContext) {
	var (
		Controller    string
		Action        string
		Serial        string
		Hash          string
		ok            bool
		internalError string
	)

	requestContext := &types.RequestContext{
		Socket:        context.Socket,
		Base:          context.Base,
		EmployeeToken: context.EmployeeToken,
		Updater:       sm.Updater,
		Env:           context.Env,
	}

	defer func() {
		if len(internalError) > 0 {
			requestContext.Response.Controller = "SocketManager"
			requestContext.Response.Action = "do"
			requestContext.Response.Error = internalError
		}

		sm.WriteMessage(context, &requestContext.Response)
	}()

	receivedJSON := context.ReceivedJSON.(map[string]interface{})

	if Serial, ok = receivedJSON["Serial"].(string); !ok {
		sm.log("do: json doesn't have action field")
		internalError = "json doesn't have `Serial` field"
		return
	}
	requestContext.Response.Serial = Serial

	if Controller, ok = receivedJSON["Controller"].(string); !ok {
		sm.log("do: json doesn't have controller field")
		internalError = "json doesn't have `controller` field"
		return
	}
	requestContext.Response.Controller = Controller

	if Action, ok = receivedJSON["Action"].(string); !ok {
		sm.log("do: json doesn't have action field")
		internalError = "json doesn't have `action` field"
		return
	}
	requestContext.Response.Action = Action

	if requestContext.ReceivedData, ok = receivedJSON["Data"]; !ok {
		sm.log("do: json doesn't have data field")
		internalError = "json doesn't have `Data` field"
		return
	}

	if Hash, ok = receivedJSON["Hash"].(string); !ok {
		sm.log("do: json doesn't have data field")
		internalError = "json doesn't have `Hash` field"
		return
	}

	data, err := internal.InsecureMarshal(requestContext.ReceivedData)
	if err != nil {
		sm.log(err.Error())
	}

	if Hash != internal.GetHash(data) && false {
		internalError = "hash sum does not add up"
		return
	}

	if _, ok = sm.controllers[Controller]; !ok {
		sm.log("do: Controller " + Controller + " doesn't exist")
		internalError = "Controller " + Controller + " doesn't exist"
		return
	}

	if _, ok = sm.controllers[Controller][Action]; !ok {
		sm.log("do: Action " + Action + " doesn't exist")
		internalError = "Action " + Action + " doesn't exist"
		return
	}

	context.LastController = Controller
	sm.controllers[Controller][Action](requestContext)
}

func (sm *SocketManager) duHast() {
}

func (sm *SocketManager) duHastMich() {

}

func (sm *SocketManager) webSocketWorker(context *types.WebSocketContext) {
	err := context.Socket.WriteJSON(struct {
		Text string
	}{
		Text: GetTongueTwister(),
	})
	if err != nil {
		sm.log("webSocketWorker" + err.Error())
		return
	}

	session := &SocketSession{
		Context: context,
	}
	sm.Sessions = append(sm.Sessions, session)

	for {
		err := context.Socket.ReadJSON(&context.ReceivedJSON)
		if err != nil {
			sm.log("webSocketWorker" + err.Error())
			return
		}

		sm.do(context)
	}
}

func (sm *SocketManager) webSocketUpdater(context *types.WebSocketContext) {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(context.Gin.Writer, context.Gin.Request, nil)
	if err != nil {
		sm.log("webSocketUpdater " + err.Error())
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			sm.log(err.Error())
		}
	}()

	context.Socket = conn
	sm.webSocketWorker(context)
}
