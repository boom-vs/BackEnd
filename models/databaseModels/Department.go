package databaseModels

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Name       string `gorm:"not null"`
	Active     uint
	SortNumber uint       `gorm:"not null; default 200"`
	Employees  []Employee `json:"-"`
}

func (d *Department) GetActive() bool {
	return d.Active == 2
}

func (d *Department) SetActive(status bool) {
	if status {
		d.Active = 2
		return
	}
	d.Active = 1
}

func (d *Department) BeforeUpdate(db *gorm.DB) (err error) {
	return nil
}
