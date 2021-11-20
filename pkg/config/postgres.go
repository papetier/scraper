package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Database        string
	Host            string
	Port            uint
	Password        string
	RefreshInterval string
	User            string
}

var DB *DBConfig

func (d *DBConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.Database)
}

func loadDbConfig() {
	DB = &DBConfig{
		Database: viper.GetString("POSTGRES_DATABASE"),
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetUint("POSTGRES_PORT"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		User:     viper.GetString("POSTGRES_USER"),
	}
}
