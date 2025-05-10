package main

import (
	"github.com/dimasmir03/web-calculator-server/internal/application"
)

//	@title			Web Calculator API
//	@version		1.0
//	@description	Distributed arithmetic expressions calculation system

//	@contact.name	API Support
//	@contact.url	http://localhost:8080
//	@contact.email	support@calculator.ru

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes	http
func main() {
	webCalculator := application.New()
	webCalculator.Run()
}
