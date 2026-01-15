package db

import (
	"context"
	"server/config"
	"server/utils"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var MilvusClient *milvusclient.Client

func initMilvusDB(cfg *config.Config) (*milvusclient.Client, error) {

	dsn := cfg.GetMilvusDSN()
	utils.Logger.Infof("use Milvus DSN:%s", dsn)

	clientConfig := &milvusclient.ClientConfig{
		Address: dsn,
	}
	client, err := milvusclient.New(context.Background(), clientConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
