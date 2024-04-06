package main

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/httpServer"
	"log"
)

func main() {

	viperInstance, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot load config. Error: {%s}", err.Error())
	}

	cfg, err := config.ParseConfig(viperInstance)
	if err != nil {
		log.Fatalf("Cannot parse config. Error: {%s}", err.Error())
	}

	s := httpServer.NewServer(cfg)
	if err = s.Run(); err != nil {
		log.Print(err)
	}
}
