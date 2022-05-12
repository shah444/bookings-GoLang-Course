package handlers

import (
	"net/http"

	"github.com/shah444/bookings/pkg/config"
	"github.com/shah444/bookings/pkg/models"
	"github.com/shah444/bookings/pkg/render"
)

// Repository is the respository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home handles the requests on the / endpoint.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About handles the requests on the /about endpoint.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remoteIP"] = remoteIP
	stringMap["test"] = "Hello World"
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{StringMap: stringMap})
}