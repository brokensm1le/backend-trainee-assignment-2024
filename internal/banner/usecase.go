package banner

type Usecase interface {
	GetBanner(params *GetBannerParams) (*string, error)
	CreateBanner(params *CreateBannerParams) (int64, error)
}
