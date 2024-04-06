package storagePostgres

import (
	"backend-trainee-assignment-2024/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

// ------------------------------------------------------------------------------------------------------------------------------

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SSLMode)

	return sqlx.Connect(c.Postgres.PgDriver, connectionUrl)
}

func CreateTable(db *sqlx.DB) error {
	var (
		query = `
		CREATE TABLE IF NOT EXISTS "user"
		(
			id         serial       not null unique,
			username   varchar(100) not null unique,
			password   varchar(100) not null,
			first_name varchar(100) not null,
			last_name  varchar(100) not null,
			email      varchar(100),
			phone      varchar(100)
		);`
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------------
