package databaseModels

import "gorm.io/gorm"

type DepartmentReportGroup struct {
	gorm.Model
	Name        string
	Departments []Department `gorm:"many2many:department_report_group_department;"`
}

func (drp *DepartmentReportGroup) BeforeUpdate(db *gorm.DB) (err error) {

	return nil
}
