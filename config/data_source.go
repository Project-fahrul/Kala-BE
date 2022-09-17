package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gorm_db *gorm.DB = nil

func create() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("SQL_HOST"), os.Getenv("SQL_PORT"), os.Getenv("SQL_USER"), os.Getenv("SQL_PASSWORD"), os.Getenv("SQL_DATABASENAME"))

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "kala.",
			SingularTable: false,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	sql, err := db.DB()

	if err != nil {
		panic(err)
	}

	sql.SetMaxIdleConns(5)
	sql.SetMaxOpenConns(5)
	sql.SetConnMaxIdleTime(5 * time.Minute)

	gorm_db = db
}

func DataSource_New() *gorm.DB {
	if gorm_db == nil {
		create()
	}
	return gorm_db
}
