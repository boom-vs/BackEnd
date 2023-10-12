package databaseModels

import (
	"crypto/sha1"
	"fmt"
	"gorm.io/gorm"
	"io"
	"time"
)

type File struct {
	gorm.Model
	Name       string
	Size       int64
	Token      string
	HashSum    string
	RealSize   int64
	UploaderId uint
	Data       []byte
}

func (f *File) getToken() {
	f.Token = ""

	hashHandler := sha1.New()
	_, err := io.WriteString(hashHandler, time.Now().String())
	if err != nil {
		return
	}

	f.Token = fmt.Sprintf("%x", hashHandler.Sum(nil))
}

func (f *File) BeforeCreate(_ *gorm.DB) (err error) {
	f.getToken()
	return
}
