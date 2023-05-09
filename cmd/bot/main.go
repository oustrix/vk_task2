package main

import (
	"log"
	"vk_task2/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error while reading config: %v", err)
	}
}
