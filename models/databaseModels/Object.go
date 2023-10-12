package databaseModels

import "gorm.io/gorm"

type ObjectGroup struct {
	gorm.Model
	Objects []Object `gorm:"foreignkey:ID"`
}

type Object struct {
	gorm.Model
	Name       string
	Active     uint8 `gorm:"not null"`
	SortNumber uint  `gorm:"not null; default 200"`
	//ID         uint  `gorm:"primarykey"`
	Icon string
}

func (o *Object) SetActive(state bool) {
	if state {
		o.Active = 2
		return
	}
	o.Active = 1
}

func (o *Object) GetActive() bool {
	return o.Active == 2
}
