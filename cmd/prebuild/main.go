package main

import (
	"crm-backend/internal"
)

func main() {

	models := internal.GetModels()
	internal.AutoControllers(models)
	models = internal.GetModels()
	internal.AutoGetDb(models)
	internal.AutoRoute(models)
	return
}
