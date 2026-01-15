package db

import (
	"os"
	"server/config"
	"server/service"
	"server/utils"
)

type InitDBHandler struct{}

func InitAllDB() {
	handler := &InitDBHandler{}
	handler.initPgSql()
	handler.initRedis()
	// handler.InitMilvus()
}

func (this *InitDBHandler) initPgSql() {
	var err error
	if PgSqlDB, err = connectPgSqlDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to PostgreSQL:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to PostgreSQL")

	if err = service.InitializeDBTables(PgSqlDB); err != nil {
		utils.Logger.Errorf("Failed to create PostgreSQL tables:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to create PostgreSQL tables")
}
func (this *InitDBHandler) initRedis() {
	var err error
	if RedisClient, err = connectRedisDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Redis:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Redis")

}
func (this *InitDBHandler) InitMilvus() {
	var err error
	if MilvusClient, err = connectMilvusDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Milvus:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Milvus")
}
