package banner

type Repository interface {
	GetAllBanners() (*[]GetFilteredBannersResponse, error)
	GetContentBanner(params *GetBannerParams) (*string, error)
	GetContentBannerAdmin(params *GetBannerParams) (*string, error)
	GetBanner(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetBannerAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersTID(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersTIDAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersFID(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	GetFilteredBannersFIDAdmin(params *GetFilteredBannersParams) (*[]GetFilteredBannersResponse, error)
	CreateBanner(params *CreateBannerParams) (int64, error)
	DeleteBanner(id int64) error
	UpdateBanner(params *UpdateBannerParams) error

	GetVersionByID(id int64) (*[]byte, error)
	GetVersionByBannerID(id int64) (*[]GetVersionResponse, error)
	DeleteVersion(versionId int64) error
}
