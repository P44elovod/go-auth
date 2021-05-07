package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/p44elovod/auth-with-gopg/pkg/handler"
	"github.com/p44elovod/auth-with-gopg/pkg/repository"
	"github.com/p44elovod/auth-with-gopg/pkg/service"
	"github.com/p44elovod/auth-with-gopg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env vriables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	rc, err := repository.NewRedisClient(repository.RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       viper.GetViper().GetInt("redis.db"),
	})

	repos := repository.NewRepository(db, rc)
	services := service.NewService(repos)
	handlers := handler.NewHandler(ctx, services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("smth went wrong when server started: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()

}
