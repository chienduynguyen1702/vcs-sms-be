package seed

import (
	"log"

	"gorm.io/gorm"
)

func InitData(db *gorm.DB) {
	var err error

	// Create roles
	err = db.Create(&Roles).Error
	if err != nil {
		log.Println("Error creating roles")
		panic(err)
	}
	log.Println("Roles created")

	// Create organizations
	err = db.Create(&Org).Error
	if err != nil {
		log.Println("Error creating organizations")
		panic(err)
	}

	// Create users
	err = db.Create(&Users).Error
	if err != nil {
		log.Println("Error creating users")
		panic(err)
	}
	log.Println("Users created")

	// Create servers
	err = db.Create(&Servers).Error
	if err != nil {
		log.Println("Error creating servers")
		panic(err)
	}
	log.Println("Servers created")
	
}
