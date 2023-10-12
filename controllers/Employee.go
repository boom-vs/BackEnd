package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"github.com/mitchellh/mapstructure"
	"log"
)

type ControllerEmployee struct {
}

func (ce *ControllerEmployee) log(message string) {
	log.Fatalln(message)
}

func (ce *ControllerEmployee) setOne(context *types.RequestContext,
	employee *webSocketModels.Employee) (*webSocketModels.Employee, error) {
	employeeGormModel := &databaseModels.Employee{}
	internal.Copier(employee, employeeGormModel)
	employeeGormModel.SetActive(employee.Active)

	tx := internal.GormUpdateOrCreate(context.Base, employeeGormModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return employee, nil
}

func (ce *ControllerEmployee) Set(context *types.RequestContext) {
	employees := []webSocketModels.Employee{}

	err := internal.MapToStrcut(context.ReceivedData, &employees)

	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	var updater []interface{}

	for _, employee := range employees {
		rEmployee, err := ce.setOne(context, &employee)
		if err != nil {
			context.Response.Error = err.Error()
			return
		}
		if rEmployee == nil {
			continue
		}
		rEmployee.Password = ""
		updater = append(updater, rEmployee)
	}
	context.Updater("Employee", updater)
}

func (ce *ControllerEmployee) GetList(context *types.RequestContext) {
	var employees []databaseModels.Employee

	context.Base.Order("sort_number desc").Find(&employees)
	for employeesIndex := 0; employeesIndex < len(employees); employeesIndex++ {
		employees[employeesIndex].Password = ""
	}

	for _, employee := range employees {
		response := &webSocketModels.Employee{}
		internal.Copier(&employee, response)
		response.Active = employee.GetActive()
		context.Response.Data = append(context.Response.Data, response)
	}
}

func (ce *ControllerEmployee) Remove(context *types.RequestContext) {
	var employees []struct{ ID uint }

	mapToStructConfig := &mapstructure.DecoderConfig{
		ErrorUnset: true,
		Result:     &employees,
	}

	mapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)
	if err != nil {
		ce.log("Remove: " + err.Error())
	}
	err = mapToStruct.Decode(context.ReceivedData)

	for _, employee := range employees {
		context.Base.Delete(&databaseModels.Employee{}, employee.ID)
	}
}
