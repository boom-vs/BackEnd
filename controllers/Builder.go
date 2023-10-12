package controllers

import (
	"crm-backend/models/databaseModels"
	"crm-backend/types"
	"gorm.io/gorm/clause"
	"log"
)

type ControllerBuilder struct {
}

func (cb *ControllerBuilder) log(message string) {
	log.Println(message)
}

func (cb *ControllerBuilder) getSubSections(context *types.RequestContext, sections []*databaseModels.Section) error {
	for sectionIndex := range sections {
		err := context.Base.Model(sections[sectionIndex]).Order("sort_number asc, title asc").
			Association("Sections").Find(&sections[sectionIndex].Sections, "active = 2")
		if err != nil {
			return err
		}

		err = cb.getSubSections(context, sections[sectionIndex].Sections)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cb *ControllerBuilder) getSections(context *types.RequestContext) ([]databaseModels.Section, error) {
	var sections []databaseModels.Section
	err := context.Base.Order("sort_number asc, title asc").Preload(clause.Associations).
		Find(&sections, "parent_id IS NULL and active = 2").Error
	if err != nil {
		return nil, err
	}

	for sectionIndex := 0; sectionIndex < len(sections); sectionIndex++ {
		err = cb.getSubSections(context, sections[sectionIndex].Sections)
		if err != nil {
			return sections, err
		}
	}

	return sections, nil
}

func (cb *ControllerBuilder) GetSections(context *types.RequestContext) {
	var (
		err      error
		sections []databaseModels.Section
	)
	sections, err = cb.getSections(context)
	if err != nil {
		context.Response.Error = "No any sections"
		log.Println(err)
		return
	}

	if len(sections) == 0 {
		context.Response.Error = "No any sections"
		return
	}

	for _, section := range sections {
		context.Response.Data = append(context.Response.Data, section)
	}
}
