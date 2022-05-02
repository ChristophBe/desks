package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var connection *gorm.DB

func initConnection() (db *gorm.DB, err error) {

	host := util.GetStringEnvironmentVariable("DB_HOST", "localhost")
	username := util.GetStringEnvironmentVariable("DB_USERNAME", "postgres")
	password := util.GetStringEnvironmentVariable("DB_PASSWORD", "password123")
	databaseName := util.GetStringEnvironmentVariable("DB_NAME", "desks_api")
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s sslmode=disable", host, username, password, databaseName)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return
}

func GetConnection(ctx context.Context) (_ *gorm.DB, err error) {
	if connection == nil {
		connection, err = initConnection()
		if err != nil {
			return
		}

	}
	return connection.WithContext(ctx), nil
}

func InitDatabase() error {
	db, err := GetConnection(context.Background())
	if err != nil {
		return err
	}

	user := new(models.User)
	room := new(models.Room)
	booking := new(models.Booking)

	return db.AutoMigrate(user, room, booking)
}
