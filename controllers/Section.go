package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"gorm.io/gorm/clause"
	"log"
)

type ControllerSection struct {
}

func (cs *ControllerSection) log(message string) {
	log.Println(message)
}

func (cs *ControllerSection) setOne(context *types.RequestContext,
	section *webSocketModels.Section) (*webSocketModels.Section, error) {
	result := &databaseModels.Section{}

	err := internal.Copier(section, result)
	if err != nil {
		return nil, err
	}

	tx := internal.GormUpdateOrCreate(context.Base, &result)
	if tx.Error != nil {
		return nil, tx.Error
	}

	err = internal.Copier(result, section)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (cs *ControllerSection) Set(context *types.RequestContext) {
	var request []webSocketModels.Section

	err := internal.MapToStrcut(context.ReceivedData, &request)
	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	var updater []interface{}

	for _, section := range request {
		item, err := cs.setOne(context, &section)
		if err != nil {
			context.Response.Error = err.Error()
			return
		}
		updater = append(updater, item)
	}
	context.Updater("Section", updater)
}

func (cs *ControllerSection) getSubSections(context *types.RequestContext, sections []*databaseModels.Section) error {
	for sectionIndex := range sections {
		err := context.Base.Model(sections[sectionIndex]).Order("sort_number desc, title desc").
			Association("Sections").Find(&sections[sectionIndex].Sections, "active = 2")
		if err != nil {
			return err
		}

		err = cs.getSubSections(context, sections[sectionIndex].Sections)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *ControllerSection) getSections(context *types.RequestContext) ([]*databaseModels.Section, error) {
	var sections []*databaseModels.Section
	err := context.Base.Order("sort_number desc, title desc").Preload(clause.Associations).
		Find(&sections, "parent_id IS NULL and active = 2").Error
	if err != nil {
		return nil, err
	}

	for sectionIndex := 0; sectionIndex < len(sections); sectionIndex++ {
		err = cs.getSubSections(context, sections[sectionIndex].Sections)
		if err != nil {
			return sections, err
		}
	}

	return sections, nil
}

func (cs *ControllerSection) GetList(context *types.RequestContext) {

	sections, err := cs.getSections(context)
	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	for _, section := range sections {
		response := &webSocketModels.Section{}
		err := internal.Copier(section, response)
		if err != nil {
			cs.log(err.Error())
		}
		context.Response.Data = append(context.Response.Data, response)
	}
}
