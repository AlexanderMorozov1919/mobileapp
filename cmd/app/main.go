// @title ClinicHub API
// @version 1.0.0
// @description API для работы с приёмами пациентов
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	_ "github.com/AlexanderMorozov1919/mobileapp/docs"
	_ "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/handlers"
	"github.com/AlexanderMorozov1919/mobileapp/internal/app"
)

func main() {
	app.New().Run()
}
