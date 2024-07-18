package db

import (
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb(host string, port string, user string, pwd string, dbName string) (*gorm.DB, error) {
	sslmode := "disable"
    timeZone := "Africa/Lagos"
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s Timezone=%s", user, pwd, dbName, host, port, sslmode, timeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
        return nil, err
    }
	log.Println("Database Connected Successfully!")
	return db, nil

}