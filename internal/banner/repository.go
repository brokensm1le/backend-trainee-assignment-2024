package banner

type Repository interface {
	GetContentBanner(params *GetBannerParams) (*string, error)
	GetContentBannerAdmin(params *GetBannerParams) (*string, error)
	CreateBanner(params *CreateBannerParams) (int64, error)
	GetBanner(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetBannerAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersTID(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersTIDAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersFID(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersFIDAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	DeleteBanner(id int64) error
	UpdateUser(params *UpdateBannerParams) error
}
