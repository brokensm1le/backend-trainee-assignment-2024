package banner

type Usecase interface {
	GetBanner(params *GetBannerParams) (*string, error)
	CreateBanner(params *CreateBannerParams) (int64, error)
	GetFilteredBanners(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	DeleteBanner(id int64) error
	UpdateUser(params *UpdateBannerParams) error
}
