package databaseModels

import "gorm.io/gorm"

//????

type TypesOfAccruals struct {
	gorm.Model
	Name        string
	Active      uint8
	SortNumber  uint
	ManualID    uint
	Icon        string
	RevenueName string
	ReportName  string
}

func (toa *TypesOfAccruals) GetActive() bool {
	return toa.Active == 2
}

func (toa *TypesOfAccruals) SetActive(state bool) {
	if state {
		toa.Active = 2
		return
	}
	toa.Active = 1
}

type TypesOfPayment struct {
	Name       string
	Active     uint8
	SortNumber uint
	ManualID   uint
	Icon       string
}
