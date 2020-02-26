package app

import (
	"bdbu/models"
)

// Copy start database migration from source to destination
func Copy(tableName string, dropTables bool) {
	if !Config.IsConfigured() {
		return
	}

	d := NewDuration()
	defer d.Completed("Completed in %v\n")

	db := models.Connect(Config.Source, Config.Destination)
	defer db.Close()

	db.Migrate(tableName, dropTables)
}
