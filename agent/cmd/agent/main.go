package main

import (
	"log"

	"github.com/dimasmir03/web-calculator-agent/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
