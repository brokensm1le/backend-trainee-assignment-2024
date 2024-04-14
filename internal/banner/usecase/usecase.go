package usecase

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/internal/banner/cache"
	"backend-trainee-assignment-2024/internal/cconstant"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type BannerUsecase struct {
	repo  banner.Repository
	cache cache.Cache
}

func NewBannerUsecase(repo banner.Repository, cache cache.Cache) banner.Usecase {
	return &BannerUsecase{repo: repo, cache: cache}
}

func (u *BannerUsecase) GetBanner(params *banner.GetBannerParams) (*string, error) {
	if !params.UseLastRevision {
		getBanner, err := u.cache.GetBanner(params.FeatureID, params.TagID)
		if err == nil {
			return &getBanner.Content, nil
		}
	}

	if params.Role == cconstant.RoleAdmin {
		return u.repo.GetContentBannerAdmin(params)
	}
	return u.repo.GetContentBanner(params)
}

func (u *BannerUsecase) GetFilteredBanners(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	if params.FeatureID != -1 && params.TagID != -1 {
		if !params.UseLastRevision {
			getBanner, err := u.cache.GetBanner(params.FeatureID, params.TagID)
			if err == nil {
				return &[]banner.GetFilteredBannersResponse{getBanner}, nil
			}
		}

		if params.Role == cconstant.RoleAdmin {
			return u.repo.GetBannerAdmin(params)
		}
		return u.repo.GetBanner(params)
	} else if params.FeatureID == -1 && params.TagID != -1 {
		if !params.UseLastRevision {
			getBanner, err := u.cache.GetBannersByTID(params.TagID)
			if err == nil {
				return &getBanner, nil
			}
		}

		if params.Role == cconstant.RoleAdmin {
			return u.repo.GetFilteredBannersTIDAdmin(params)
		}
		return u.repo.GetFilteredBannersTID(params)
	}

	if !params.UseLastRevision {
		getBanner, err := u.cache.GetBannersByFID(params.FeatureID)
		if err == nil {
			return &getBanner, nil
		}
	}
	if params.Role == cconstant.RoleAdmin {
		return u.repo.GetFilteredBannersFIDAdmin(params)
	}
	return u.repo.GetFilteredBannersFID(params)
}

func (u *BannerUsecase) CreateBanner(params *banner.CreateBannerParams) (int64, error) {
	return u.repo.CreateBanner(params)
}

func (u *BannerUsecase) DeleteBanner(id int64) error {
	return u.repo.DeleteBanner(id)
}

func (u *BannerUsecase) UpdateUser(params *banner.UpdateBannerParams) error {
	return u.repo.UpdateBanner(params)
}

func (u *BannerUsecase) GetVersionByBannerID(id int64) (*[]banner.GetVersionDecodeResponse, error) {
	dataVersions, err := u.repo.GetVersionByBannerID(id)
	if err != nil {
		return nil, err
	}
	var res []banner.GetVersionDecodeResponse
	for _, version := range *dataVersions {
		var data banner.GetVersionDecodeResponse
		err = json.Unmarshal(version.Data, &data)
		log.Println(data)
		if err != nil {
			return nil, err
		}
		data.VersionID = version.VersionID
		res = append(res, data)
	}
	return &res, nil
}

func (u *BannerUsecase) DeleteVersion(versionId int64) error {
	return u.repo.DeleteVersion(versionId)
}

func (u *BannerUsecase) SelectVersion(versionId int64, bannerId int64) error {
	dataBannerBytes, err := u.repo.GetVersionByID(versionId)
	if err != nil {
		return err
	}
	var data banner.GetFilteredBannersResponse
	err = json.Unmarshal(*dataBannerBytes, &data)
	if err != nil {
		return err
	}

	tagsStr := strings.Split(strings.Trim(data.TagIDs, "{}"), ",")
	tags := make([]int64, len(tagsStr))
	for i, s := range tagsStr {
		tags[i], _ = strconv.ParseInt(s, 10, 64)
	}

	err = u.repo.DeleteVersion(versionId)
	if err != nil {
		return err
	}

	err = u.repo.UpdateBanner(&banner.UpdateBannerParams{
		BannerID:  bannerId,
		TagIDs:    tags,
		FeatureID: float64(data.FeatureID),
		Content:   data.Content,
		IsActive:  data.IsActive,
	})
	if err != nil {
		return err
	}

	return nil
}
