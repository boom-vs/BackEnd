package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"fmt"
	"log"
)

type ControllerAuthorization struct {
}

func (ca *ControllerAuthorization) log(message string) {
	log.Println(message)
}

func (ca *ControllerAuthorization) Get(context *types.RequestContext) {
	var getRequest []webSocketModels.AuthorizationGetRequest

	err := internal.MapToStrcut(context.ReceivedData, &getRequest)

	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	if len(getRequest) != 1 {
		context.Response.Error = "Wrong input data" + fmt.Sprintf("%d", len(getRequest))
		return
	}

	session := databaseModels.Session{}
	context.Base.Preload("Employee").First(&session, "token = ?",
		getRequest[0].Token)

	if session.ID == 0 {
		context.Response.Error = fmt.Sprintf("The session with token `%s` was not found", getRequest[0].Token)
		return
	}

	if !session.GetActive() {
		context.Response.Error = fmt.Sprintf("The session with token `%s` isn't active", getRequest[0].Token)
		return
	}

	context.EmployeeToken = getRequest[0].Token
	response := &webSocketModels.AuthorizationGetResponse{
		ID: session.ID,
	}

	internal.Copier(&session.Employee, response)
	context.Response.Data = append(context.Response.Data, response)
	return
}
