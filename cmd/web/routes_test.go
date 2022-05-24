package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shah444/bookings-GoLang-Course/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
		case *chi.Mux:
			// Do nothing. Test passed.
		default:
			t.Error(fmt.Sprintf("type is not *chi.Mux, it is %T", v))
	}
}