package repository

import (
	"backend-trainee-assignment-2024/internal/banner"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) banner.Repository {
	return &postgresRepository{db: db}
}
