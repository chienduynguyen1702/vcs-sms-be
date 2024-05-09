package configs

import (
	"log"

	"github.com/chienduynguyen1702/vcs-sms-be/models"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {

	err := db.AutoMigrate(&models.Organization{})
	if err != nil {
		log.Println("Failed to migrate Organization models")
		return err
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("Failed to migrate User models")
		return err
	}
	log.Printf("Database migrated\n")
	return nil
}
