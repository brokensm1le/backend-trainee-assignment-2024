package banner

import (
	"time"
)

type GetBannerParams struct {
	Token           string `json:"-"`
	TagID           int64  `query:"tag_id" db:"tag_id"`
	FeatureID       int64  `query:"feature_id" db:"feature_id"`
	UseLastRevision bool   `query:"use_last_revision"`
	Role            int    `json:"-"`
}

type GetFilteredBannersParams struct {
	Token           string `json:"-"`
	TagID           int64  `query:"tag_id" db:"tag_id"`
	FeatureID       int64  `query:"feature_id" db:"feature_id"`
	Limit           int64  `query:"limit"`
	Offset          int64  `query:"offset"`
	UseLastRevision bool   `query:"use_last_revision"`
	Role            int    `json:"-"`
}

type GetFilteredBannersResponse struct {
	BannerID  int64     `json:"banner_id" db:"banner_id"`
	TagIDs    string    `json:"tag_ids" db:"tag_ids"`
	FeatureID int64     `json:"feature_id" db:"feature_id"`
	Content   string    `json:"content" db:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GetFilteredBannersDecodeResponse struct {
	BannerID  int64     `json:"banner_id" db:"banner_id"`
	TagIDs    []int64   `json:"tag_ids" db:"tag_ids"`
	FeatureID int64     `json:"feature_id" db:"feature_id"`
	Content   string    `json:"content" db:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateBannerParams struct {
	Token     string      `json:"-"`
	BannerID  int64       `json:"banner_id" db:"banner_id"`
	TagIDs    interface{} `json:"tag_ids" db:"tag_ids"`
	FeatureID interface{} `json:"feature_id" db:"feature_id"`
	Content   interface{} `json:"content" db:"content"`
	IsActive  interface{} `json:"is_active" db:"is_active"`
	UpdatedAt time.Time   `json:"-" db:"updated_at"`
}

type CreateBannerParams struct {
	Token           string  `json:"-"`
	TagIDs          []int64 `json:"tag_ids"`
	FeatureID       int64   `json:"feature_id" db:"feature_id"`
	Content         string  `json:"content" db:"content"`
	IsActive        bool    `json:"is_active" db:"is_active"`
	UseLastRevision bool    `json:"use_last_revision"`
}

type DBBanner struct {
	BannerID  int64     `json:"banner_id" db:"banner_id"`
	TagIDs    string    `json:"tag_ids" db:"tag_ids"`
	FeatureID int64     `json:"feature_id" db:"feature_id"`
	Content   string    `json:"content" db:"content"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GetVersionResponse struct {
	VersionID int64  `json:"version_id" db:"version_id"`
	Data      []byte `json:"data" db:"data"`
}

type SelectVersionParams struct {
	BannerID  int64 `query:"banner_id"`
	VersionID int64 `query:"version_id"`
}

type GetVersionDecodeResponse struct {
	VersionID int64     `json:"version_id"`
	BannerID  int64     `json:"banner_id"`
	TagIDs    string    `json:"tag_ids"`
	FeatureID int64     `json:"feature_id"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
