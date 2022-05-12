package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shah444/bookings/pkg/config"
	"github.com/shah444/bookings/pkg/handlers"
	"github.com/shah444/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var session *scs.SessionManager

var app config.AppConfig
// main is the main application fn
func main() {

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
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}