package dbrepo

import (
	"database/sql"

	"github.com/shah444/bookings-GoLang-Course/internal/config"
	"github.com/shah444/bookings-GoLang-Course/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}