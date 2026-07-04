package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/iokreon1/coffee-pos/config"
	"github.com/iokreon1/coffee-pos/internal/handler"
	"github.com/iokreon1/coffee-pos/pkg/database"
	redispkg "github.com/iokreon1/coffee-pos/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewMySQL(cfg.MysqlDSN())
	if err != nil {
		log.Fatalf("Failed to connect MySQL: %v", err)
	}
	defer db.Close()

	rdb, err := redispkg.NewRedis(cfg.RedisAddr(), cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	defer rdb.Close()

	r := handler.NewRouter(cfg.AppEnv)

	log.Printf("Starting server on port :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
