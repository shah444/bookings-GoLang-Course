package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shah444/bookings-GoLang-Course/internal/config"
	"github.com/shah444/bookings-GoLang-Course/internal/driver"
	"github.com/shah444/bookings-GoLang-Course/internal/handlers"
	"github.com/shah444/bookings-GoLang-Course/internal/helpers"
	"github.com/shah444/bookings-GoLang-Course/internal/models"
	"github.com/shah444/bookings-GoLang-Course/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var session *scs.SessionManager

var app config.AppConfig

// main is the main application fn
func main() {
	db, err := run()
	if err != nil {
		log.Fatal("Error running application")
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// What I am going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Change this to true during production
	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database
	log.Println("connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=viditshah password=")
	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}

	log.Println("connected to database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
