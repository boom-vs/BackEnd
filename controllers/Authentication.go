package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"fmt"
)

type ControllerAuthentication struct {
}

func (ca *ControllerAuthentication) log(message string) {
	fmt.Println(message)
}

func (ca *ControllerAuthentication) Login(context *types.RequestContext) {
	var request = []webSocketModels.AuthenticationLoginRequest{}

	err := internal.MapToStrcut(context.ReceivedData, &request)
	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	if len(request) != 1 {
		context.Response.Error = "Wrong input data" + fmt.Sprintf("%d", len(request))
		return
	}

	var employee databaseModels.Employee
	tx := context.Base.First(&employee, "login = ?", request[0].Login)
	if tx.Error != nil {
		context.Response.Error = tx.Error.Error()
		return
	}

	if !employee.GetActive() {
		context.Response.Error = "User inactive"
		return
	}

	if !employee.CheckPassowrdWithHash(request[0].Password) {
		context.Response.Error = fmt.Sprintf("Wrong password for `%s` user", request[0].Login)
		return
	}

	session := &databaseModels.Session{
		Employee: employee,
	}
	session.SetKeep(request[0].Keep)
	session.SetActive(true)
	db := context.Base.Create(session)

	if session.ID == 0 {
		context.Response.Error = db.Error.Error()
		return
	}

	context.Response.Data = append(context.Response.Data,
		webSocketModels.AuthenticationLoginResponce{Token: session.Token})
	context.EmployeeToken = session.Token
	context.EmployeeId = employee.ID
}

func (ca *ControllerAuthentication) Logout(context *types.RequestContext) {
	session := &databaseModels.Session{}
	context.Base.First(session, "token = ?", context.EmployeeToken)

	if session.ID != 0 {
		session.SetActive(false)
		context.Base.Save(session)

	}
	context.EmployeeToken = ""
}
