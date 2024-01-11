package config

import "fmt"

type Config struct {
	ServerConfig ServerConfig `mapstructure:"server"`
	DBConfig     DBConfig     `mapstructure:"db"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DBConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

func (dbConf DBConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DBName)
}
