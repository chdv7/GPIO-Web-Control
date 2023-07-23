package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chdv7/GPIO-Web-Control/internal/config"
	"github.com/chdv7/GPIO-Web-Control/internal/gpio"
	"github.com/go-chi/chi/v5"
)

var Gpio gpio.Gpio
var Config map[string]config.GPIO

type keyType string

const gpioPathKey keyType = "gpioPath"
const valueKey keyType = "value"

func SetupRoutes(r chi.Router) {
	r.Route("/gpio", func(r chi.Router) {
		r.Route("/{gpioID}", func(r chi.Router) {
			r.Use(gpioCtx)
			r.Get("/", getGPIO)
			r.Put("/", setGPIO)
		})
	})
}

func gpioCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gpioID := chi.URLParam(r, "gpioID")

		if _, ok := Config[gpioID]; !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), gpioPathKey, Config[gpioID])
		value := r.URL.Query().Get("value")
		if value != "" {
			ctx = context.WithValue(ctx, valueKey, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getGPIO(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gpioPath, ok := ctx.Value(gpioPathKey).(config.GPIO)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	bytes, err := Gpio.ReadGPIO(gpioPath.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func setGPIO(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gpioPath, ok := ctx.Value(gpioPathKey).(config.GPIO)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	value, ok := ctx.Value(valueKey).(string)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	if value != "0" && value != "1" {
		http.Error(w, "Invalid value", http.StatusInternalServerError)
		return
	}

	err := Gpio.WriteGPIO(gpioPath.Location, []byte(value))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}

	bytes, err := Gpio.ReadGPIO(gpioPath.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func GpioConfigToMap(cfg config.GPIOConfig) map[string]config.GPIO {
	gpioMap := make(map[string]config.GPIO)
	for _, gpio := range cfg.GPIOs {
		gpioMap[gpio.Name] = gpio
	}
	return gpioMap
}
