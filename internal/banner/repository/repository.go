package repository

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/customTime"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) banner.Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) GetBanner(params *banner.GetBannerParams) (*string, error) {
	var (
		data  string
		query = `
		SELECT content
		FROM %[1]s 
		WHERE feature_id = $1 AND $2 = ANY (tag_ids);
		`

		values = []any{params.FeatureID, params.TagID}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Get(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) CreateBanner(params *banner.CreateBannerParams) (int64, error) {
	var (
		query = `
		INSERT INTO %[1]s (feature_id, content, is_active, tag_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING banner_id`

		values = []any{params.FeatureID, params.Content, params.IsActive,
			"{" + strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ", ", -1), "[]") + "}",
			customTime.GetMoscowTime(), customTime.GetMoscowTime(),
		}
		id int64
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	fmt.Println(query)

	if err := p.db.Get(&id, query, values...); err != nil {
		return 0, err
	}

	return id, nil
}
