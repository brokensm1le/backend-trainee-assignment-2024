package storagePostgres

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/banner"
	"fmt"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"log"
	"math/rand"
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
		);
		CREATE TABLE IF NOT EXISTS "auth"
		(
			id		  	    	bigserial    not null unique,
			name				text		 not null,
			email		   		varchar(255) not null unique,
			password   			text		 not null,
			role      			smallint	 not null default 0,
			refresh_token    	text		 not null unique,
			refresh_token_ttl   timestamp	 not null
		);
		`
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

type PairId struct {
	FID int64
	TID int64
}

func GenerateTable(c *config.Config, repo banner.Repository) error {

	checker := make(map[PairId]bool)

	ans := make([]banner.CreateBannerParams, 0)

	for i := int64(0); i < c.Generator.CntRow; i++ {
		ban := banner.CreateBannerParams{}
		fid := RandomInt64(c.Generator.CntFeature)
		ban.FeatureID = fid

		randomCntTag := RandomInt64(c.Generator.MaxTagInRow)
		sliceTags := make([]int64, 0)
		for j := int64(0); j <= randomCntTag; j++ {
			tid := RandomInt64(c.Generator.CntTag)
			if _, ok := checker[PairId{fid, tid}]; !ok {
				checker[PairId{fid, tid}] = true
				sliceTags = append(sliceTags, tid)
			}
		}
		if len(sliceTags) == 0 {
			continue
		}
		ban.TagIDs = sliceTags
		ban.Content = fmt.Sprintf("GIGA-Content %d", i)
		ban.IsActive = RandomBool()

		ans = append(ans, ban)
	}

	for _, p := range ans {
		_, err := repo.CreateBanner(&p)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------------------------------------

// анлак не получилось
//func convertJSONToCSV(data *[]banner.CreateBannerParams, nameOutFile string) error {
//	outputFile, err := os.Create(nameOutFile)
//	if err != nil {
//		return err
//	}
//	defer outputFile.Close()
//
//	writer := csv.NewWriter(outputFile)
//	defer writer.Flush()
//
//	header := []string{"feature_id", "content", "is_active", "created_at", "updated_at", "tag_ids"}
//	if err := writer.Write(header); err != nil {
//		return err
//	}
//
//	for _, r := range *data {
//		var csvRow []string
//		csvRow = append(csvRow, strconv.FormatInt(r.FeatureID, 10), r.Content, strconv.FormatBool(r.IsActive))
//		csvRow = append(csvRow, customTime.GetMoscowTime().String(), customTime.GetMoscowTime().String())
//		csvRow = append(csvRow, strings.ReplaceAll(fmt.Sprint(r.TagIDs), " ", ","))
//		if err := writer.Write(csvRow); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func RandomInt64(limit int64) int64 {
	return rand.Int63() % limit
}

func RandomBool() bool {
	return rand.Int()%2 == 0
}
