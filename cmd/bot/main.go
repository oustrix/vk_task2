package main

import (
	"log"
	"vk_task2/config"
	"vk_task2/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error while reading config: %v", err)
	}

	a := app.NewApp(cfg)
	a.Run()
}
