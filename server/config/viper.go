package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var VP *viper.Viper

func InitViper(FileName string) (*viper.Viper, *Config) {

	v := viper.New()
	v.SetConfigName(FileName)
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AutomaticEnv()
	fmt.Printf("Config name set to %s\n", FileName)
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error config file: %v\n", err)
		os.Exit(0)
	}
	fmt.Println("Configuration file read successfully")

	var s Config
	if err := v.Unmarshal(&s); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error unmarshal config: %v\n", err)
		os.Exit(0)
	}
	fmt.Println("Configuration unmarshaled successfully")
	return v, &s
}
