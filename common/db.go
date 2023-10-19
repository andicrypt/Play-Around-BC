package common

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func MustConnectDatabaseWithName(dbConfig *Database, dbName string, testing bool) (*gorm.DB, error) {
	var (
		err error
		db  *gorm.DB
	)
	// load sqlite db for testing purpose
	if testing {
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbConfig.Host, dbConfig.User, dbConfig.Password, dbName, dbConfig.Port)
		dialect := postgres.Open(dsn)
		db, err = gorm.Open(dialect, &gorm.Config{})
		if err != nil {
			panic(err)
		}
		pgDB, err := db.DB()
		if err != nil {
			panic(err)
		}

		pgDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)
		pgDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		pgDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		log.Error("error querying SELECT 1", "err", err)
		panic(err)
	}
	return db, err
}

func createPgDb(cfg *Database) {
	db, err := MustConnectDatabaseWithName(cfg, "postgres", false)
	if err != nil {
		panic(err)
	}
	if db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)).Error != nil {
		log.Error("error while creating database", "err", err, "dbName", cfg.DBName)
	}
}

func NewDBConn(dbConfig *Database, testing bool) (*gorm.DB, error) {
	// load sqlite db for testing purpose
	if testing {
		return gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	// create db
	createPgDb(dbConfig)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)
	dialect := postgres.Open(dsn)
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	pgDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	if dbConfig.ConnMaxLifetime > 0 {
		pgDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Hour)
	}
	pgDB.SetConnMaxIdleTime(2 * time.Minute)
	pgDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	pgDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	err = db.Raw("SELECT 1").Error
	if err != nil {
		log.Error("error querying SELECT 1", "err", err)
		panic(err)
	}
	return db, err
}