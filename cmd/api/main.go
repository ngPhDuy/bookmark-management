// @title Bookmark Management API
// @version 1.0
// @description REST API for managing bookmarks
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ngPhDuy/bookmark-management/internal/api"
	"github.com/ngPhDuy/bookmark-management/internal/config"
	"github.com/ngPhDuy/bookmark-management/internal/infra/redis"
)

func main() {
	cfg := config.NewAppConfig()

	ctx := context.Background()

	rdb, err := redis.NewRedisClient(ctx, &redis.RedisConfig{
		Host:     cfg.RedisHost,
		Port:     fmt.Sprintf("%d", cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err != nil {
		log.Fatal("Redis connection error: ", err)
	}

	defer rdb.Close()

	engine := api.NewEngine(cfg, rdb)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}
