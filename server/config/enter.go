package config

import "fmt"

var Settings *Config

type Config struct {
	SYSTEM_IS_DEV                bool   `mapstructure:"SYSTEM_IS_DEV"`
	SYSTEM_SERVER_PORT           int    `mapstructure:"SYSTEM_SERVER_PORT"`
	JWT_SIGNING_KEY              string `mapstructure:"JWT_SIGNING_KEY"`
	JWT_EXPIRES_TIME             string `mapstructure:"JWT_EXPIRES_TIME"`
	JWT_BUFFER_TIME              string `mapstructure:"JWT_BUFFER_TIME"`
	JWT_ISSUER                   string `mapstructure:"JWT_ISSUER"`
	REDIS_HOST                   string `mapstructure:"REDIS_HOST"`
	REDIS_PORT                   int    `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD               string `mapstructure:"REDIS_PASSWORD"`
	REDIS_DB                     int    `mapstructure:"REDIS_DB"`
	EMAIL_TO                     string `mapstructure:"EMAIL_TO"`
	EMAIL_PORT                   int    `mapstructure:"EMAIL_PORT"`
	EMAIL_FROM                   string `mapstructure:"EMAIL_FROM"`
	EMAIL_HOST                   string `mapstructure:"EMAIL_HOST"`
	EMAIL_IS_SSL                 bool   `mapstructure:"EMAIL_IS_SSL"`
	EMAIL_SECRET                 string `mapstructure:"EMAIL_SECRET"`
	EMAIL_NICKNAME               string `mapstructure:"EMAIL_NICKNAME"`
	CAPTCHA_KEY_LONG             int    `mapstructure:"CAPTCHA_KEY_LONG"`
	CAPTCHA_IMG_WIDTH            int    `mapstructure:"CAPTCHA_IMG_WIDTH"`
	CAPTCHA_IMG_HEIGHT           int    `mapstructure:"CAPTCHA_IMG_HEIGHT"`
	CAPTCHA_OPEN_CAPTCHA         int    `mapstructure:"CAPTCHA_OPEN_CAPTCHA"`
	CAPTCHA_OPEN_CAPTCHA_TIMEOUT int    `mapstructure:"CAPTCHA_OPEN_CAPTCHA_TIMEOUT"`
	POSTGRES_HOST                string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT                int    `mapstructure:"POSTGRES_PORT"`
	POSTGRES_DB                  string `mapstructure:"POSTGRES_DB"`
	POSTGRES_USERNAME            string `mapstructure:"POSTGRES_USERNAME"`
	POSTGRES_PASSWORD            string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_MAX_IDLE_CONNS      int    `mapstructure:"POSTGRES_MAX_IDLE_CONNS"`
	POSTGRES_MAX_OPEN_CONNS      int    `mapstructure:"POSTGRES_MAX_OPEN_CONNS"`
	MILVUS_HOST                  string `mapstructure:"MILVUS_HOST"`
	MILVUS_PORT                  int    `mapstructure:"MILVUS_PORT"`
}

func (this *Config) GetServerPort() string {
	return fmt.Sprintf(":%d", this.SYSTEM_SERVER_PORT)
}

func (this *Config) GetPostgreDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		this.POSTGRES_HOST, this.POSTGRES_PORT, this.POSTGRES_USERNAME, this.POSTGRES_PASSWORD, this.POSTGRES_DB,
	)
}

func (this *Config) GetRedisDSN() string {
	return fmt.Sprintf("%s:%d", this.REDIS_HOST, this.REDIS_PORT)
}

func (this *Config) GetMilvusDSN() string {
	return fmt.Sprintf("%s:%d", this.MILVUS_HOST, this.MILVUS_PORT)
}
