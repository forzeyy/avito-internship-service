package main

import (
	"github.com/forzeyy/avito-internship-service/internal/app"
	"github.com/forzeyy/avito-internship-service/internal/config"
)

func main() {
	app.Run(config.LoadConfig())
}
