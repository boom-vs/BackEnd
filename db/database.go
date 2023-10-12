package db

import (
	"crm-backend/models/databaseModels"
	"crm-backend/types"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DataBase struct {
	Handler *gorm.DB
	Config  *gorm.Config
}

func (db *DataBase) getDsn(ctx *types.GormContext) string {

	//Secure option. if ssl disabled then only local db
	if !ctx.SSLMode {
		ctx.Host = "localhost"
	}

	sslmode := "disable"
	if ctx.SSLMode {
		sslmode = "enable"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		ctx.Host,
		ctx.User,
		ctx.Password,
		ctx.Database,
		ctx.Port,
		sslmode,
	)
}

func (db *DataBase) log(message string) {
	log.Println(message)
}

func (db *DataBase) initData() error {

	//Sections
	var sections []databaseModels.Section
	result := db.Handler.Find(&sections)

	if result.Error != nil {
		return result.Error
	}

	//Insert sections
	if len(sections) == 0 {
		newSections := db.getSections()
		db.Handler.Create(&newSections)
	}

	//Employee
	var employee []databaseModels.Employee
	result = db.Handler.Find(&employee)

	if result.Error != nil {
		return result.Error
	}

	if len(employee) == 0 {
		newEmployee := databaseModels.Employee{
			FirstName:      "Иван",
			LastName:       "Иванов",
			PatronymicName: "Иванович",
			Login:          "admin",
			Password:       "admin",
			Position:       "",
			Department: databaseModels.Department{
				Name: "admin",
			},
		}
		db.Handler.Create(&newEmployee)
	}

	return nil
}

func (db *DataBase) Connect(ctx *types.GormContext) {
	var err error
	db.Handler, err = gorm.Open(postgres.Open(db.getDsn(ctx)), db.Config)
	if err != nil {
		db.Handler = nil
		return
	}

	if ctx.AutoMigrate {
		err = db.AutoMigrate()
		if err != nil {
			db.log(err.Error())
			return
		}

		err = db.initData()
		if err != nil {
			db.log(err.Error())
			return
		}
	}
}
