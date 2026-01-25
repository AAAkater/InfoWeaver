package db

import (
	"os"
	"server/config"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

type InitDBHandler struct{}

func InitAllDB() {
	handler := &InitDBHandler{}
	handler.initPgSql()
	handler.initRedis()
	handler.initMinio()
	// handler.InitMilvus()
}

func (this *InitDBHandler) initPgSql() {
	var err error
	if PgSqlDB, err = connectPgSqlDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to PostgreSQL:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to PostgreSQL")

	// Auto-migrate database schema
	if err = PgSqlDB.AutoMigrate(&models.User{}, &models.File{}, &models.Document{}, &models.Memory{}); err != nil {
		utils.Logger.Errorf("Failed to create PostgreSQL tables:%s", err)
		os.Exit(0)
	}

	// Create initial admin account if not exists
	var admin models.User
	result := PgSqlDB.Where("role = ?", "admin").First(&admin)
	if result.Error != gorm.ErrRecordNotFound {
		return
	}

	// Create admin user with default credentials
	res := PgSqlDB.Create(&models.User{
		Username: config.Settings.SYSTEM_ADMIN_NAME,
		Email:    config.Settings.SYSTEM_ADMIN_EMAIL,
		Password: utils.BcryptHash(config.Settings.SYSTEM_ADMIN_PASSWORD),
		Role:     "admin",
	})

	if res.Error != nil {
		utils.Logger.Errorf("Failed to create admin account:%s", res.Error)
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
func (this *InitDBHandler) initMinio() {
	var err error
	if MinioClient, err = connectMinioClient(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Minio:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Minio")
}
func (this *InitDBHandler) InitMilvus() {
	var err error
	if MilvusClient, err = connectMilvusDB(config.Settings); err != nil {
		utils.Logger.Errorf("Failed to connect to Milvus:%s", err)
		os.Exit(0)
	}
	utils.Logger.Info("success to connect to Milvus")
}
