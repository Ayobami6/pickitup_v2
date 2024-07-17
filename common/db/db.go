package db

import (
	"fmt"
	"log"

	"github.com/Ayobami6/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "ayo")
	pwd := config.GetEnv("DB_PWD", "password")
	dbName := config.GetEnv("DB_NAME", "pickitup_db")
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