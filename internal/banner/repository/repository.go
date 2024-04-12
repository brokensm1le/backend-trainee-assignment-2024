package repository

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/internal/cconstant"
	"backend-trainee-assignment-2024/pkg/customTime"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"strings"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) banner.Repository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) GetAllBanners() (*[]banner.GetFilteredBannersResponse, error) {
	var (
		data  []banner.GetFilteredBannersResponse
		query = `
		SELECT *
		FROM %[1]s 
		WHERE is_active = true;
		`
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Select(&data, query); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetContentBanner(params *banner.GetBannerParams) (*string, error) {
	var (
		data  string
		query = `
		SELECT content
		FROM %[1]s 
		WHERE feature_id = $1 AND $2 = ANY (tag_ids) AND is_active = true;
		`

		values = []any{params.FeatureID, params.TagID}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Get(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetContentBannerAdmin(params *banner.GetBannerParams) (*string, error) {
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
		WHERE feature_id = $1 AND $2 = ANY (tag_ids) AND is_active = true;
		`

		values = []any{params.FeatureID, params.TagID}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Select(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}

func (p *postgresRepository) GetBannerAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
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
		WHERE feature_id = $1 AND is_active = true
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

func (p *postgresRepository) GetFilteredBannersFIDAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
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
		log.Println("ERROR(3 ADMIN):", err)
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
		WHERE $1 = ANY (tag_ids) AND is_active = true
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

func (p *postgresRepository) GetFilteredBannersTIDAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
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
		query = `INSERT INTO %[1]s (feature_id, content, is_active, tag_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING banner_id;`

		values = []any{params.FeatureID, params.Content, params.IsActive,
			"{" + strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ", ", -1), "[]") + "}",
			customTime.GetMoscowTime(), customTime.GetMoscowTime(),
		}
		id int64
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	if _, err := tx.Exec(fmt.Sprintf("LOCK TABLE %[1]s IN ROW EXCLUSIVE MODE;", cconstant.BannerDB)); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.QueryRow(query, values...).Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	queryGet := `SELECT count(*) from banner
                        WHERE EXISTS (SELECT UNNEST(tag_ids) INTERSECT SELECT UNNEST($1::bigint[])) 
                          AND feature_id = $2
                        ;`
	valuesGet := []any{pq.Array(params.TagIDs), params.FeatureID}

	var cntRow int
	if err := tx.QueryRow(queryGet, valuesGet...).Scan(&cntRow); err != nil {
		tx.Rollback()
		return 0, err
	}

	if cntRow != 1 {
		tx.Rollback()
		return 0, fmt.Errorf("you already have that relationship")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
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

	chngBanner, err := p.getBannerById(params.BannerID)

	if params.TagIDs != nil {
		nameParams = nameParams + "tag_ids"
		cntParams++
		chngBanner.TagIDs = "{" + strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ", ", -1), "[]") + "}"
		values = append(values, "{"+strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ", ", -1), "[]")+"}")
	}
	if params.FeatureID != nil {
		if cntParams > 0 {
			nameParams = nameParams + ", feature_id"
		}
		cntParams++
		chngBanner.FeatureID = int64(params.FeatureID.(float64))
		values = append(values, int64(params.FeatureID.(float64)))
	}
	if params.Content != nil {
		if cntParams > 0 {
			nameParams = nameParams + ", content"
		}
		cntParams++
		values = append(values, params.Content)
	}
	if params.IsActive != nil {
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

	// -----------------------------------------------------------------------------------------------------------------------------

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(fmt.Sprintf("LOCK TABLE %[1]s IN ROW EXCLUSIVE MODE;", cconstant.BannerDB)); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(query, values...); err != nil {
		tx.Rollback()
		return err
	}

	queryGet := `SELECT count(*) from banner
                        WHERE EXISTS (SELECT UNNEST(tag_ids) INTERSECT SELECT UNNEST($1::bigint[])) 
                          AND feature_id = $2
                        ;`
	valuesGet := []any{chngBanner.TagIDs, chngBanner.FeatureID}

	var cntRow int
	if err := tx.QueryRow(queryGet, valuesGet...).Scan(&cntRow); err != nil {
		tx.Rollback()
		return err
	}
	if cntRow != 1 {
		tx.Rollback()
		return fmt.Errorf("you already have that relationship")
	}

	return tx.Commit()
}

func (p *postgresRepository) getBannerById(id int64) (*banner.GetFilteredBannersResponse, error) {
	var (
		data  banner.GetFilteredBannersResponse
		query = `
		SELECT *
		FROM %[1]s 
		WHERE banner_id = $1;
		`

		values = []any{id}
	)

	query = fmt.Sprintf(query, cconstant.BannerDB)

	if err := p.db.Get(&data, query, values...); err != nil {
		return &data, err
	}

	return &data, nil
}
