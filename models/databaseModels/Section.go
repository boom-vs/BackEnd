package databaseModels

import "gorm.io/gorm"

type Section struct {
	gorm.Model
	Name       string
	Title      string
	Active     uint8
	SortNumber uint
	ParentID   *uint
	Icon       string
	Sections   []*Section `gorm:"foreignkey:ParentID"`
}

func (s *Section) GetActive() bool {
	return s.Active == 2
}

func (s *Section) SetActive(state bool) {
	if state {
		s.Active = 2
		return
	}
	s.Active = 1
}
