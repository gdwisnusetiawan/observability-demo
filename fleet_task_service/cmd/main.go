package main

import (
	"fathil/go-observability/fleet_task_service/config"
	"fathil/go-observability/fleet_task_service/internal/app"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	app := app.New(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
