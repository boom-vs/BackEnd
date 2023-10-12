package databaseModels

import (
	"crypto/sha1"
	"fmt"
	"gorm.io/gorm"
	"io"
	"time"
)

type Session struct {
	gorm.Model
	EmployeeId  uint
	Employee    Employee
	Token       string    `gorm:"not null"`
	TimeCreated time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Keep        uint8     `gorm:"not null"`
	Active      uint8     `gorm:"not null"`
}

func (s *Session) SetKeep(value bool) {
	if value {
		s.Keep = 2
		return
	}
	s.Keep = 1
}

func (s *Session) GetKeep() bool {
	if s.Keep == 1 {
		return false
	}
	return true
}

func (s *Session) SetActive(value bool) {
	if value {
		s.Active = 2
		return
	}
	s.Active = 1
}

func (s *Session) GetActive() bool {
	if s.Active == 1 {
		return false
	}
	return true
}

func (s *Session) getToken() error {
	s.Token = ""

	hashHandler := sha1.New()
	_, err := io.WriteString(hashHandler, time.Now().String())
	if err != nil {
		return err
	}

	s.Token = fmt.Sprintf("%x", hashHandler.Sum(nil))
	return nil
}

func (s *Session) BeforeCreate(_ *gorm.DB) (err error) {
	if err := s.getToken(); err != nil {
		return err
	}
	return nil
}
