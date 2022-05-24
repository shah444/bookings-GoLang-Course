package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shah444/bookings-GoLang-Course/internal/config"
	"github.com/shah444/bookings-GoLang-Course/internal/handlers"
	"github.com/shah444/bookings-GoLang-Course/internal/models"
	"github.com/shah444/bookings-GoLang-Course/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var session *scs.SessionManager

var app config.AppConfig
// main is the main application fn
func main() {
	err := run()
	if err != nil {
		log.Fatal("Error running application")
	}

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
		// What I am going to put in the session
		gob.Register(models.Reservation{})
	
		// Change this to true during production
		app.InProduction = false
	
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
	
		app.Session = session
	
		tc, err := render.CreateTemplateCache()
		if err != nil {
			log.Fatal("Cannot create template cache")
			return err
		}
	
		app.TemplateCache = tc
		app.UseCache = false
	
		repo := handlers.NewRepo(&app)
		handlers.NewHandlers(repo)
	
		render.NewTemplate(&app)

		return nil
}