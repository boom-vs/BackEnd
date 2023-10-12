package databaseModels

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	//Employee Employee
	Data string
}
