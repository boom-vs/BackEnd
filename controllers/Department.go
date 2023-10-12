package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
)

type ControllerDepartment struct {
}

func (cd *ControllerDepartment) log(message string) {
	log.Println(message)
}

func (cd *ControllerDepartment) setOne(context *types.RequestContext,
	department *webSocketModels.DepartmentRequest) (error, *webSocketModels.DepartmentRequest) {
	newDepartment := &databaseModels.Department{}

	err := internal.Copier(department, newDepartment)
	if err != nil {
		return err, nil
	}

	tx := internal.GormUpdateOrCreate(context.Base, newDepartment)
	if tx.Error != nil {
		return tx.Error, nil
	}

	err = internal.Copier(newDepartment, department)
	return nil, department
}

func (cd *ControllerDepartment) Set(context *types.RequestContext) {
	var request []webSocketModels.DepartmentRequest

	mapToStructConfig := &mapstructure.DecoderConfig{
		ErrorUnset: true,
		Result:     &request,
	}

	mapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)
	if err != nil {
		cd.log(err.Error())
	}

	err = mapToStruct.Decode(context.ReceivedData)

	if err != nil {
		context.Response.Error = fmt.Sprintf("Invalid structure: `%+v` %s", context.ReceivedData, err.Error())
		return
	}

	var updater []interface{}

	for _, department := range request {
		err, rDepartment := cd.setOne(context, &department)
		if err != nil {
			context.Response.Error = err.Error()
			return
		}
		updater = append(updater, rDepartment)
	}
	context.Updater("Department", updater)
}

func (cd *ControllerDepartment) GetList(context *types.RequestContext) {
	var departments []*databaseModels.Department

	context.Base.Order("sort_number, name").Find(&departments)
	for _, department := range departments {
		departmentResponse := &webSocketModels.DepartmentResponse{}
		internal.Copier(department, departmentResponse)
		departmentResponse.ID = department.ID
		context.Response.Data = append(context.Response.Data, departmentResponse)
	}
}
