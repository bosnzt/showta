package model

import (
	// "fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	syslog "log"
	"os"
	"path/filepath"
	"showta.cc/app/lib/util"
	"showta.cc/app/system/conf"
	"time"
)

var db *gorm.DB

func InitDb(cfg conf.Database) {
	fullDbName := conf.AbsPath(cfg.Dbname)
	checkDbDir(fullDbName)

	newLogger := logger.New(
		syslog.New(os.Stdout, "\r\n", syslog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	var err error
	db, err = gorm.Open(sqlite.Open(fullDbName), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	db.Exec("PRAGMA journal_mode=WAL;")
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Storage{}, &FolderSetting{}, &Preference{})
}

func checkDbDir(pathStr string) {
	dirName, _ := filepath.Split(pathStr)
	if dirName == "" {
		return
	}

	ok, _ := util.PathExist(dirName)
	if ok {
		return
	}

	os.MkdirAll(dirName, os.ModePerm)
}
