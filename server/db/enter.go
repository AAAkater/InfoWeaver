package db

import (
	"os"
	"server/config"
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
	if PgSqlDB, err = initPgSqlDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to PostgreSQL:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to PostgreSQL")
}
func (this *InitDBHandler) initRedis() {
	var err error
	if RedisClient, err = initRedisDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Redis:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Redis")

}
func (this *InitDBHandler) InitMilvus() {
	var err error
	if MilvusClient, err = initMilvusDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Milvus:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Milvus")
}
