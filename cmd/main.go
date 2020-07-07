package main

import (
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phuwn/tools/db"
	"github.com/phuwn/tools/log"

	"github.com/newzoom/api/pkg/handler"
	"github.com/newzoom/api/pkg/handler/ws"
)

func main() {
	db.Start()
	defer db.Close()

	addr := ":8080"
	log.Status("listening on port%s", addr)
	go ws.NewHub()

	err := http.ListenAndServe(addr, handler.Router())
	if err != nil {
		log.Status("server got terminated, err: %s", err.Error())
	}
}
