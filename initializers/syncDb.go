package initializers

import (
	"github.com/FerranMarin/generic-dating/models"
)

func SyncDb() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Swipe{})
	DB.AutoMigrate(&models.Match{})
}
