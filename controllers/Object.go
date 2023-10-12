package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"log"
)

type ControllerObject struct {
}

func (co *ControllerObject) log(message string) {
	log.Println(message)
}

func (co *ControllerObject) setOne(context *types.RequestContext,
	object *webSocketModels.Object) (error, *webSocketModels.Object) {
	result := &databaseModels.Object{}

	internal.Copier(object, result)
	result.SetActive(object.Active)

	tx := internal.GormUpdateOrCreate(context.Base, result)

	if tx.Error != nil {
		return tx.Error, nil
	}
	err := internal.Copier(result, object)
	if err != nil {
		return err, nil
	}
	return nil, object
}

func (co *ControllerObject) Set(context *types.RequestContext) {
	var request []webSocketModels.Object

	err := internal.MapToStrcut(context.ReceivedData, &request)

	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	var updater []interface{}

	for _, object := range request {

		err, item := co.setOne(context, &object)
		if err != nil {
			context.Response.Error = err.Error()
			return
		}
		updater = append(updater, item)
	}
	context.Updater("Object", updater)
}

func (co *ControllerObject) GetList(context *types.RequestContext) {
	var objects []*databaseModels.Object

	context.Base.Order("sort_number, name").Find(&objects)
	for _, object := range objects {
		response := &webSocketModels.Object{}
		internal.Copier(object, response)
		response.Active = object.GetActive()
		context.Response.Data = append(context.Response.Data, response)
	}
}
