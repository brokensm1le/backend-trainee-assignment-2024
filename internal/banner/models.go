package banner

import "time"

type GetBannerParams struct {
	Token           string `json:"-"`
	TagID           int64  `query:"tag_id" db:"tag_id"`
	FeatureID       int64  `query:"feature_id" db:"feature_id"`
	UseLastRevision bool   `query:"use_last_revision"`
	Role            int    `json:"-"`
}

type GetFilteredBannersParams struct {
	Token     string `json:"-"`
	TagID     int64  `query:"tag_id" db:"tag_id"`
	FeatureID int64  `query:"feature_id" db:"feature_id"`
	Limit     int64  `query:"limit"`
	Offset    int64  `query:"offset"`
	Role      int    `json:"-"`
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
