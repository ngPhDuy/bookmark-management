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
