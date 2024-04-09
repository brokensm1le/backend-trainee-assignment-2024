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

func (p *postgresRepository) GetContentBanner(params *banner.GetBannerParams) (*string, error) {
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

func (p *postgresRepository) GetBanner(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	var (
		data  []banner.GetFilteredBannersResponse
		query = `
		SELECT *
		FROM %[1]s 
		WHERE feature_id = $1 AND $2 = ANY (tag_ids);
		`

		values = []any{params.FeatureID, params.TagID}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetFilteredBannersFID(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	var (
		data  []banner.GetFilteredBannersResponse
		query = `
		SELECT *
		FROM %[1]s 
		WHERE feature_id = $1
		LIMIT $2
		OFFSET $3;
		`

		values = []any{params.FeatureID, params.Limit, params.Offset}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetFilteredBannersTID(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	var (
		data  []banner.GetFilteredBannersResponse
		query = `
		SELECT *
		FROM %[1]s 
		WHERE $1 = ANY (tag_ids)
		LIMIT $2
		OFFSET $3;
		`

		values = []any{params.TagID, params.Limit, params.Offset}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Select(&data, query, values...); err != nil {
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

func (p *postgresRepository) DeleteBanner(id int64) error {
	var (
		query = `
		DELETE FROM %[1]s 
		WHERE banner_id = $1
		`

		values = []any{id}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) UpdateUser(params *banner.UpdateBannerParams) error {
	var (
		query string = `
		UPDATE %[1]s SET
		`

		cntParams  int
		nameParams string
		values     []any
	)

	fmt.Println("params:", params)
	tagIds, ok := params.TagIDs.([]int64)
	fmt.Println(tagIds, ok)
	if params.TagIDs != nil {
		fmt.Println("Check1")
		nameParams = nameParams + "tag_ids"
		cntParams++
		values = append(values, "{"+strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ", ", -1), "[]")+"}")
	}
	if params.FeatureID != nil {
		fmt.Println("Check2")
		if cntParams > 0 {
			nameParams = nameParams + ", feature_id"
		}
		cntParams++
		values = append(values, int64(params.FeatureID.(float64)))
	}
	if params.Content != nil {
		fmt.Println("Check3")
		if cntParams > 0 {
			nameParams = nameParams + ", content"
		}
		cntParams++
		values = append(values, params.Content)
	}
	if params.IsActive != nil {
		fmt.Println("Check4")
		if cntParams > 0 {
			nameParams = nameParams + ", is_active"
		}
		cntParams++
		values = append(values, params.IsActive)
	}
	if cntParams > 0 {
		nameParams = nameParams + ", updated_at"
		cntParams++
		values = append(values, customTime.GetMoscowTime())
	}
	fmt.Println("VALUES:", values)
	fmt.Println("namePARAMS:", nameParams)

	if cntParams > 1 {
		query += "(" + nameParams + ") = \n\t\t\t("
		query += fmt.Sprintf("$%d", 1)
		for i := 2; i <= cntParams; i++ {
			query += fmt.Sprintf(",$%d", i)
		}
		query += ")\n"
	} else {
		query += nameParams + " = \n\t\t\t"
		query += fmt.Sprintf("$%d", 1) + "\n"
	}

	values = append(values, params.BannerID)
	query += fmt.Sprintf("\t\tWHERE banner_id = $%d", len(values))
	// -----------------------------------------------------------------------------------------------------------------------------

	query = fmt.Sprintf(query, cconstant.BannerDB)

	fmt.Println("QUERY:", query)
	// -----------------------------------------------------------------------------------------------------------------------------

	if _, err := p.db.Exec(query, values...); err != nil {
		return err
	}

	return nil
}
