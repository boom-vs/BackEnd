package databaseModels

import (
	"crm-backend/types"
	"crypto/sha1"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"strings"
	"unicode"
)

type Employee struct {
	gorm.Model
	ID             uint       `gorm:"primarykey"`
	DepartmentID   uint       `gorm:""`
	ManualID       uint16     `gorm:"not null;unique"`
	LastName       string     `gorm:"not null"`
	FirstName      string     `gorm:"not null"`
	PatronymicName string     `gorm:"not null"`
	Active         uint       `gorm:"not null"`
	SortNumber     uint       `gorm:"not null; default 200"`
	Department     Department `json:"-"`
	Position       string
	PhoneNumber    string
	Telegram       string
	Email          string
	Login          string `gorm:"not null,unique"`
	Password       string `gorm:"not null"`
}

func (e *Employee) checkDepartmentID(db *gorm.DB) error {
	if e.DepartmentID == 0 {
		return errors.New("field DepartmentID must not 0")
	}

	if e.Department.ID == 0 {
		department := Department{}
		db.First(&department, "id = ?", e.DepartmentID)
		if department.ID == 0 {
			return errors.New(fmt.Sprintf("field Department can't be applied, Department with ID %d not found",
				e.DepartmentID))
		}
	}

	return nil
}

func (e *Employee) checkManualID() error {
	if e.ManualID == 0 {
		return errors.New("field ManualID must not 0")
	}
	return nil
}

func (e *Employee) checkNameSymbols(name string) error {
	badSymbolIndex := strings.IndexFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	if badSymbolIndex >= 0 {
		return errors.New(string(name[badSymbolIndex]))
	}
	return nil
}

func (e *Employee) checkLastName() error {
	if len(e.LastName) == 0 {
		return errors.New("field LastName must not have length 0")
	}

	if err := e.checkNameSymbols(e.LastName); err != nil {
		return errors.New(fmt.Sprintf("field LastName must have only letters (`%s`)", err.Error()))
	}

	return nil
}

func (e *Employee) checkFirstName() error {
	if len(e.FirstName) == 0 {
		return errors.New("field FirstName must not have length 0")
	}

	if err := e.checkNameSymbols(e.FirstName); err != nil {
		return errors.New(fmt.Sprintf("field FirstName must have only letters (`%s`)", err.Error()))
	}

	return nil
}

func (e *Employee) checkPatronymicName() error {
	if len(e.PatronymicName) == 0 {
		return errors.New("field PatronymicName must not have length 0")
	}

	if err := e.checkNameSymbols(e.PatronymicName); err != nil {
		return errors.New(fmt.Sprintf("field PatronymicName must have only letters (`%s`)", err.Error()))
	}

	return nil
}

func (e *Employee) checkActive() error {
	if e.Active == 0 || e.Active > 2 {
		return errors.New(fmt.Sprintf("field Active must not have undefined state (%d)", e.Active))
	}
	return nil
}

func (e *Employee) checkSortNumber() error {
	if e.SortNumber == 0 {
		e.SortNumber = types.DefaultSortNumber
	}
	return nil
}

func (e *Employee) checkDepartment() error {

	if e.Department.ID != 0 && e.Department.ID != e.DepartmentID {
		return errors.New(fmt.Sprintf("field Department.ID has conflict with DepartmentID"))
	}

	return nil
}

func (e *Employee) checkPosition() error {
	return nil
}

func (e *Employee) checkPhoneNumber() error {
	if len(e.PhoneNumber) < 5 {
		return errors.New("field PhoneNumber must not have length < 5")
	}

	replacer := strings.NewReplacer("-", "", "+", "", " ", "")
	e.PhoneNumber = replacer.Replace(e.PhoneNumber)

	if len(e.PhoneNumber) > 11 {
		if e.PhoneNumber[0] == '8' {
			e.PhoneNumber = "7" + e.PhoneNumber[1:]
		}

		e.PhoneNumber = "+" + e.PhoneNumber
	}

	return nil
}

func (e *Employee) checkTelegram() error {

	if len(e.Telegram) == 0 {
		return nil
	}

	if len(e.Telegram) < 5 {
		return errors.New("field Telegram must not have length < 5")
	}

	e.Telegram = strings.ReplaceAll(strings.ToLower(e.Telegram), "@", "")

	badSymbolIndex := strings.IndexFunc(e.Telegram, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || (r == '_'))
	})
	if badSymbolIndex >= 0 {
		return errors.New(fmt.Sprintf("field Telegram must have only letters (`%s`)",
			string(e.Telegram[badSymbolIndex])))
	}

	e.Telegram = "@" + e.Telegram

	return nil
}

func (e *Employee) checkEmail() error {
	if len(e.Email) < 6 {
		return errors.New("field Email must not have length < 6")
	}

	emailSplit := strings.Split(e.Email, "@")
	if len(emailSplit) != 2 {
		return errors.New("field Email must have only one symbol at (@)")
	}

	if badIndex := strings.IndexFunc(e.Email, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || (r == '-') || (r == '.') || (r == '@'))
	}); badIndex > -1 {
		return errors.New(fmt.Sprintf("field Email must not have spec.symbols `%s`", string(e.Email[badIndex])))
	}

	return nil
}

func (e *Employee) checkLogin() error {
	if len(e.Login) == 0 {
		return errors.New("field Login must not have length 0")
	}

	e.Login = strings.ToLower(e.Login)

	if badIndex := strings.IndexFunc(e.Login, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || (r == '-') || (r == '.'))
	}); badIndex > -1 {
		return errors.New(fmt.Sprintf("field Login must not have spec.symbols `%s`", string(e.Login[badIndex])))
	}

	return nil
}

func (e *Employee) checkPassword() error {
	if len(e.Password) == 0 {
		//return errors.New("field Password must not have length 0")
		return nil
	}

	if e.Password[0] == '\t' && len(e.Password) != 41 {
		return errors.New("field Password is corrupted")
	}

	if e.Password[0] != '\t' {
		hash, err := e.getHash(e.Password)
		if err != nil {
			return errors.New("field Password has error, err:" + err.Error())
		}
		e.Password = "\t" + hash
	}
	return nil
}

func (e *Employee) getHash(clearText string) (string, error) {

	if len(e.Login) == 0 {
		return "", errors.New("login is empty")
	}

	randInstance := rand.New(rand.NewSource(int64(e.Login[0])))

	salt := ""
	for i := 0; i < len(e.Login); i++ {
		salt += string(rune(int32('a') + int32(randInstance.Intn(20))))
	}

	hashHandler := sha1.New()
	_, err := io.WriteString(hashHandler, clearText+salt)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hashHandler.Sum(nil)), nil
}

func (e *Employee) CheckPassowrdWithHash(clearText string) bool {
	hash, err := e.getHash(clearText)
	if err != nil {
		return false
	}

	if hash != e.Password[1:] {
		return false
	}
	return true
}

func (e *Employee) HashIt() error {
	hash, err := e.getHash(e.Password)
	if err != nil {
		return err
	}
	e.Password = hash
	return nil
}

func (e *Employee) checkAllField(db *gorm.DB) error {
	if err := e.checkDepartmentID(db); err != nil {
		return err
	}
	if err := e.checkManualID(); err != nil {
		return err
	}
	if err := e.checkLastName(); err != nil {
		return err
	}
	if err := e.checkFirstName(); err != nil {
		return err
	}
	if err := e.checkPatronymicName(); err != nil {
		return err
	}
	if err := e.checkActive(); err != nil {
		return err
	}
	if err := e.checkSortNumber(); err != nil {
		return err
	}
	if err := e.checkDepartment(); err != nil {
		return err
	}
	if err := e.checkPosition(); err != nil {
		return err
	}
	if err := e.checkPhoneNumber(); err != nil {
		return err
	}
	if err := e.checkTelegram(); err != nil {
		return err
	}
	if err := e.checkEmail(); err != nil {
		return err
	}
	if err := e.checkLogin(); err != nil {
		return err
	}
	if err := e.checkPassword(); err != nil {
		return err
	}
	return nil
}

func (e *Employee) SetActive(status bool) {
	if status {
		e.Active = 2
		return
	}
	e.Active = 1
}

func (e *Employee) GetActive() bool {
	return e.Active == 2
}

func (e *Employee) BeforeUpdate(db *gorm.DB) (err error) {
	if e.ID == 0 {
		return errors.New("field ID must not 0")
	}

	if err := e.checkAllField(db); err != nil {
		return err
	}
	return nil
}

func (e *Employee) BeforeCreate(db *gorm.DB) (err error) {
	if err := e.checkAllField(db); err != nil {
		return err
	}

	return nil
}
