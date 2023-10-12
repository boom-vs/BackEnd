package main

import (
	"crm-backend/config"
	db2 "crm-backend/db"
	"crm-backend/webserver"
	"gorm.io/gorm"
	"log"
)

func main() {
	configReader := &config.ConfigReader{}

	db := &db2.DataBase{
		Config: &gorm.Config{},
	}
	db.Connect(configReader.GetGormContext())

	if db.Handler == nil {
		log.Println("not connected DB")
		return
	}

	webServer := &webserver.WebServer{}

	webServerContext := configReader.GetWebServerContext()
	webServerContext.DataBase = db.Handler
	webServerContext.Env = configReader.GetEnv()

	log.Println("main DB", db.Handler)

	webServer.Start(webServerContext)
}
