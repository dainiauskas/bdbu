package app

import (
	"bdbu/models"
)

func Copy(tableName string) {
	if !Config.IsConfigured() {
		return
	}

	d := NewDuration()
	defer d.Completed("Completed in %v\n")

	db := models.Connect(Config.Source, Config.Destination)
	defer db.Close()

	db.Migrate(tableName)
}
