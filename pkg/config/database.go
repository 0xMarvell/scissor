package config

import (
	"log"

	"github.com/0xMarvell/scissor/pkg/models"
	"github.com/0xMarvell/scissor/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes a database connection for the API
func Connect() {
	var dbConnectErr error

	// dsn := "host=db user=postgres password=postgres dbname=scissor port=5432 sslmode=disable TimeZone=Africa/Lagos" // for development
	dsn := "postgres://scissor_db_service_user:XZzvqtu77yCIMQwfJl3i9UBOhvgsFEmx@dpg-cicm7n59aq03rjn42to0-a.oregon-postgres.render.com/scissor_db_service" // for deployment

	DB, dbConnectErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.CheckDbErr("error connecting to database: ", dbConnectErr)
	log.Println("Database Connection Successful 🚀")
}

// RunMigrations runs migrations for the database.
func RunMigrations() {
	migrationErr := DB.AutoMigrate(&models.User{}, &models.URL{})
	utils.CheckDbErr("Migration failed: ", migrationErr)
	log.Println("Migration Successful 🚀")
}
