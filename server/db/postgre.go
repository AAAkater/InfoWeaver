package db

import (
	"server/config"
	"server/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PgSqlDB *gorm.DB

func connectPgSqlDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetPostgreDSN()
	utils.Logger.Infof("use PostgreSQL DSN:%s", dsn)

	postgreConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(postgreConfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(cfg.POSTGRES_MAX_IDLE_CONNS)
		sqlDB.SetMaxOpenConns(cfg.POSTGRES_MAX_OPEN_CONNS)
	}
	return db, nil
}
