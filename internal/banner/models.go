package banner

type GetBannerParams struct {
	Token           string `json:"-"`
	TagID           int64  `query:"tag_id" db:"tag_id"`
	FeatureID       int64  `query:"feature_id" db:"feature_id"`
	UseLastRevision bool   `query:"use_last_revision"`
}

type CreateBannerParams struct {
	Token           string  `json:"-"`
	TagIDs          []int64 `json:"tag_ids"`
	FeatureID       int64   `json:"feature_id" db:"feature_id"`
	Content         string  `json:"content" db:"content"`
	IsActive        bool    `json:"is_active" db:"is_active"`
	UseLastRevision bool    `json:"use_last_revision"`
}
