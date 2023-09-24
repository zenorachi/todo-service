package main

import (
	"github.com/zenorachi/todo-service/internal/app"
	"github.com/zenorachi/todo-service/internal/config"
)

func main() {
	app.Run(config.New())
}
