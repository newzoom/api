package main

import (
	"github.com/joho/godotenv"

	"github.com/phuwn/tools/log"
	"github.com/phuwn/tools/util"

	"github.com/newzoom/api/pkg/server"
	"github.com/newzoom/api/pkg/service"
	"github.com/newzoom/api/pkg/store"
)

// init server stuff
func init() {
	env := util.Getenv("RUN_MODE", "")
	if env == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	store := store.New()
	service := service.New()
	server.NewServerCfg(store, service)
}
