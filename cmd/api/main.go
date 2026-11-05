package main

import (
	"github.com/ngPhDuy/bookmark-management/internal/api"
	"github.com/ngPhDuy/bookmark-management/internal/config"
)

func main() {
	cfg, err := config.NewAppConfig()

	if err != nil {
		panic(err)
	}

	engine := api.NewEngine(cfg)

	if err := engine.Start(); err != nil {
		panic(err)
	}
}
