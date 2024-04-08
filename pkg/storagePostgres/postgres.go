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
		CREATE TABLE IF NOT EXISTS "banner"
		(
			banner_id       bigserial    not null unique,
			feature_id   	bigint       not null,
			tag_ids   		bigint[]	 not null,
			content 		text		 not null,
			is_active  		boolean 	 not null default true,
			created_at      timestamp	 not null,
			updated_at      timestamp	 not null
		);`
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------------
