package main

import (
	"filmLibraryVk/api/handler"
	"filmLibraryVk/internal/repository"
	"filmLibraryVk/internal/service"
	"filmLibraryVk/pkg"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can not initialize configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("can not load env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:         viper.GetString("db.host"),
		Port:         viper.GetString("db.port"),
		Username:     viper.GetString("db.username"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:       viper.GetString("db.dbname"),
		SSLMode:      viper.GetString("db.sslmode"),
		MigrationURL: viper.GetString("db.migration_url"),
	})
	if err != nil {
		log.Fatalf("can not initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(pkg.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("can not run http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}