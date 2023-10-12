package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"fmt"
	"gorm.io/gorm/clause"
	"log"
)

type ControllerDepartmentReportGroup struct {
}

func (cdrg *ControllerDepartmentReportGroup) log(message string) {
	log.Println(message)
}

func (cdrg *ControllerDepartmentReportGroup) setOne(context *types.RequestContext,
	department *webSocketModels.DepartmentReportGroup) (error, *webSocketModels.DepartmentReportGroup) {
	departmentGroup := &databaseModels.DepartmentReportGroup{}

	err := internal.Copier(department, departmentGroup)
	if err != nil {
		return err, nil
	}

	if departmentGroup.ID != 0 {
		tx := context.Base.Exec("DELETE FROM department_report_group_department WHERE "+
			"department_report_group_id = ? AND  department_id NOT IN ? ", department.ID, department.DepartmentsId)
		if tx.Error != nil {
			return tx.Error, nil
		}
	}

	for _, ids := range department.DepartmentsId {
		subDepartment := databaseModels.Department{}
		tx := context.Base.First(&subDepartment, ids)

		if tx.Error != nil {
			return tx.Error, nil
		}

		if subDepartment.ID > 0 {
			departmentGroup.Departments = append(departmentGroup.Departments, subDepartment)
		}
	}

	tx := internal.GormUpdateOrCreate(context.Base, departmentGroup)
	if tx.Error != nil {
		return tx.Error, nil
	}
	return nil, department
}

func (cdrg *ControllerDepartmentReportGroup) Set(context *types.RequestContext) {
	var request []webSocketModels.DepartmentReportGroup

	err := internal.MapToStrcut(context.ReceivedData, &request)

	if err != nil {
		context.Response.Error = err.Error()
		return
	}

	var updater []interface{}

	for _, departmentGroup := range request {
		err, rDepartment := cdrg.setOne(context, &departmentGroup)
		if err != nil {
			context.Response.Error = err.Error()
			return
		}
		updater = append(updater, rDepartment)
	}
	context.Updater("DepartmentReportGroup", updater)
}

func (cdrg *ControllerDepartmentReportGroup) GetList(context *types.RequestContext) {
	var (
		departmentReportGroup []databaseModels.DepartmentReportGroup
	)
	context.Base.Preload(clause.Associations).Find(&departmentReportGroup)

	fmt.Printf("%#v\n", departmentReportGroup)

	for _, reportGroup := range departmentReportGroup {
		departmentResponse := &webSocketModels.DepartmentReportGroup{
			ID:   reportGroup.ID,
			Name: reportGroup.Name,
		}

		for _, departments := range reportGroup.Departments {
			departmentResponse.DepartmentsId = append(departmentResponse.DepartmentsId, departments.ID)
		}

		context.Response.Data = append(context.Response.Data, departmentResponse)
	}
}
