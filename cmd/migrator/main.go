package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"os"
	_ "sso/migrations"
)

type Config struct {
	host     string
	port     int
	dbName   string
	user     string
	password string
	sslMode  string
	isDrop   bool
}

func main() {
	if err := initConfig(); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := initDbConfig()

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbName, cfg.sslMode))
	if err != nil {
		panic(err)
	}

	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if migrationsPath == "" {
		panic("storagePath and migrationsPath are required")
	}

	if !cfg.isDrop {
		if err := goose.Up(db, migrationsPath); err != nil {
			panic(err)
		}
	} else {
		if err := goose.Down(db, migrationsPath); err != nil {
			panic(err)
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("cmd/migrator")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initDbConfig() Config {
	return Config{
		host:     viper.GetString("db.host"),
		port:     viper.GetInt("db.port"),
		dbName:   viper.GetString("db.name"),
		user:     viper.GetString("db.username"),
		password: os.Getenv("DB_PASSWORD"),
		sslMode:  viper.GetString("db.ssl_mode"),
		isDrop:   viper.GetBool("db.is_drop"),
	}
}
