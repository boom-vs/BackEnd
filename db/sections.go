package db

import (
	"crm-backend/models/databaseModels"
	"github.com/goccy/go-json"
	"os"
)

func (db *DataBase) getSections() []databaseModels.Section {
	bytes, err := os.ReadFile("db/sections.json")
	if err != nil {
		db.log(err.Error())
		return nil
	}

	var sections []databaseModels.Section

	err = json.Unmarshal(bytes, &sections)
	if err != nil {
		db.log(err.Error())
		return nil
	}

	return sections
}
