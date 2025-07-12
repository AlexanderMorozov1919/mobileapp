package main

import "github.com/AlexanderMorozov1919/mobileapp/internal/app"

func main() {
	// @title Booking-service
	// @version 1.0.0
	// @description This is REST API server for working with patient's receptions
	// @contact.name API Support
	// @BasePath /api/v1/
	app.New().Run()
}
