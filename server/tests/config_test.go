package tests

import (
	"server/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSetting(t *testing.T) {
	// Test loading config from existing yaml file
	config.VP, config.Settings = config.InitViper(config.TEST_ENV_FILENAME)

	assert.NotNil(t, config.Settings, "Config should not be nil")

	// Verify Postgres configuration
	assert.Equal(t, "localhost", config.Settings.MILVUS_HOST)
	assert.Equal(t, 5432, config.Settings.POSTGRES_PORT)
	assert.Equal(t, "postgres", config.Settings.POSTGRES_USERNAME)
	assert.Equal(t, "postgres", config.Settings.POSTGRES_PASSWORD)
	assert.Equal(t, "InfoWeaver", config.Settings.POSTGRES_DB)
	// Verify Redis configuration
	assert.Equal(t, "localhost", config.Settings.REDIS_HOST)
	assert.Equal(t, 6379, config.Settings.REDIS_PORT)
	assert.Equal(t, 0, config.Settings.REDIS_DB)
	// Verify Milvus configuration
	assert.Equal(t, "localhost", config.Settings.MILVUS_HOST)
	assert.Equal(t, 19530, config.Settings.MILVUS_PORT)

	// Verify other configuration
	assert.True(t, config.Settings.SYSTEM_IS_DEV)
	assert.Equal(t, 8080, config.Settings.SYSTEM_SERVER_PORT)
}
