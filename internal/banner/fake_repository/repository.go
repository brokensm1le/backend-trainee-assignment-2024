package fake_repository

// ----- This repo for test -----

import (
	"backend-trainee-assignment-2024/internal/banner"
	"fmt"
	"strconv"
	"strings"
)

type PairIDs struct {
	FID int64
	TID int64
}

type fakeRepository struct {
	Id        int64
	Banner    map[PairIDs]banner.GetFilteredBannersResponse
	BannerFID map[int64][]banner.GetFilteredBannersResponse
	BannerTID map[int64][]banner.GetFilteredBannersResponse
}

func NewFakeRepository() banner.Repository {
	return &fakeRepository{
		Id:        0,
		Banner:    make(map[PairIDs]banner.GetFilteredBannersResponse),
		BannerFID: make(map[int64][]banner.GetFilteredBannersResponse),
		BannerTID: make(map[int64][]banner.GetFilteredBannersResponse),
	}
}

func (p *fakeRepository) GetAllBanners() (*[]banner.GetFilteredBannersResponse, error) {
	// dont use in tests (for cache)
	return nil, nil
}

func (p *fakeRepository) GetContentBanner(params *banner.GetBannerParams) (*string, error) {
	ban, ok := p.Banner[PairIDs{params.FeatureID, params.TagID}]
	if !ok || !ban.IsActive {
		fmt.Println("OK:", ok)
		fmt.Println("BANNER:", ban)
		fmt.Println("MAP", p.Banner)
		data := ""
		return &data, fmt.Errorf("sql: no rows in result set")
	}
	return &ban.Content, nil
}

func (p *fakeRepository) GetContentBannerAdmin(params *banner.GetBannerParams) (*string, error) {
	ban, ok := p.Banner[PairIDs{params.FeatureID, params.TagID}]
	if !ok {
		data := ""
		return &data, fmt.Errorf("sql: no rows in result set")
	}
	return &ban.Content, nil
}

func (p *fakeRepository) GetBanner(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	ban, ok := p.Banner[PairIDs{params.FeatureID, params.TagID}]
	if !ok || !ban.IsActive {
		return &[]banner.GetFilteredBannersResponse{}, fmt.Errorf("sql: no rows in result set")
	}
	return &[]banner.GetFilteredBannersResponse{ban}, nil
}

func (p *fakeRepository) GetBannerAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	ban, ok := p.Banner[PairIDs{params.FeatureID, params.TagID}]
	if !ok {
		return &[]banner.GetFilteredBannersResponse{}, fmt.Errorf("sql: no rows in result set")
	}
	return &[]banner.GetFilteredBannersResponse{ban}, nil
}

func (p *fakeRepository) GetFilteredBannersTID(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	ans := make([]banner.GetFilteredBannersResponse, 0)
	banners, ok := p.BannerTID[params.TagID]
	if !ok {
		return &ans, fmt.Errorf("sql: no rows in result set")
	}
	for _, response := range banners {
		if !response.IsActive {
			continue
		}
		ans = append(ans, response)
	}

	if len(ans) == 0 {
		return &ans, fmt.Errorf("sql: no rows in result set")
	}
	return &ans, nil
}

func (p *fakeRepository) GetFilteredBannersTIDAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	banners, ok := p.BannerTID[params.TagID]
	if !ok || len(banners) == 0 {
		return &banners, fmt.Errorf("sql: no rows in result set")
	}
	return &banners, nil
}

func (p *fakeRepository) GetFilteredBannersFID(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	ans := make([]banner.GetFilteredBannersResponse, 0)
	banners, ok := p.BannerFID[params.FeatureID]
	if !ok {
		return &ans, fmt.Errorf("sql: no rows in result set")
	}
	for _, response := range banners {
		if !response.IsActive {
			continue
		}
		ans = append(ans, response)
	}

	if len(ans) == 0 {
		return &ans, fmt.Errorf("sql: no rows in result set")
	}
	return &ans, nil
}

func (p *fakeRepository) GetFilteredBannersFIDAdmin(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	banners, ok := p.BannerFID[params.FeatureID]
	if !ok || len(banners) == 0 {
		return &banners, fmt.Errorf("sql: no rows in result set")
	}
	return &banners, nil
}

func (p *fakeRepository) CreateBanner(params *banner.CreateBannerParams) (int64, error) {
	p.Id++
	ban := banner.GetFilteredBannersResponse{
		BannerID:  p.Id,
		TagIDs:    strings.Trim(strings.Replace(fmt.Sprint(params.TagIDs), " ", ",", -1), "[]"),
		FeatureID: params.FeatureID,
		Content:   params.Content,
		IsActive:  params.IsActive,
	}
	for _, tag := range params.TagIDs {
		p.Banner[PairIDs{params.FeatureID, tag}] = ban
		if _, ok := p.BannerTID[tag]; !ok {
			p.BannerTID[tag] = make([]banner.GetFilteredBannersResponse, 0)
		}
		p.BannerTID[tag] = append(p.BannerTID[tag], ban)
	}
	if _, ok := p.BannerFID[params.FeatureID]; !ok {
		p.BannerFID[params.FeatureID] = make([]banner.GetFilteredBannersResponse, 0)
	}
	p.BannerFID[params.FeatureID] = append(p.BannerFID[params.FeatureID], ban)
	fmt.Println("CREATE:", p.Banner)
	return p.Id, nil
}

func (p *fakeRepository) DeleteBanner(id int64) error {
	var fid int64
	tidsStr := ""
	tids := make([]int64, 0)
	for _, response := range p.Banner {
		if response.BannerID == id {
			fid = response.FeatureID
			tidsStr = response.TagIDs
			tagIds := strings.Split(response.TagIDs, ",")
			for _, tagId := range tagIds {
				tid, _ := strconv.ParseInt(tagId, 10, 64)
				tids = append(tids, tid)
			}
			break
		}
	}
	if tidsStr == "" {
		return fmt.Errorf("sql: no rows in result set")
	}

	fidSlice := p.BannerFID[fid]
	deleteID := -1
	for i, response := range fidSlice {
		if response.TagIDs == tidsStr {
			deleteID = i
			break
		}
	}
	fidSlice = append(fidSlice[:deleteID], fidSlice[deleteID+1:]...)
	p.BannerFID[fid] = fidSlice

	for _, tid := range tids {
		delete(p.Banner, PairIDs{fid, tid})
		tidSlice := p.BannerTID[tid]
		for i := 0; i < len(tidSlice); i++ {
			if tidSlice[i].FeatureID == fid {
				p.BannerTID[tid] = append(tidSlice[:i], tidSlice[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (p *fakeRepository) UpdateUser(params *banner.UpdateBannerParams) error {
	// without this method
	return nil
}
