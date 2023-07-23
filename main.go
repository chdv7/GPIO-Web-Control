package main

import (
	"log"
	"net/http"

	"github.com/chdv7/GPIO-Web-Control/internal/config"
	simple_gpio "github.com/chdv7/GPIO-Web-Control/internal/gpio/simple-gpio"
	"github.com/chdv7/GPIO-Web-Control/internal/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router.Config = router.GpioConfigToMap(config.GetConfig())
	router.Gpio = simple_gpio.GpioConnector{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	router.SetupRoutes(r)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
