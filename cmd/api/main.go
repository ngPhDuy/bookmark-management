// @title Bookmark Management API
// @version 1.0
// @description REST API for managing bookmarks
// @host localhost:8080
// @BasePath /
package main

import (
	"github.com/ngPhDuy/bookmark-management/internal/api"
	"github.com/ngPhDuy/bookmark-management/internal/config"
)

func main() {
	cfg := config.NewAppConfig()

	engine := api.NewEngine(cfg)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}
