package lazyCache

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/internal/banner/cache"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type PairIDs struct {
	FID int64
	TID int64
}

type Cache struct {
	mutex     sync.RWMutex
	mutexFID  sync.RWMutex
	mutexTID  sync.RWMutex
	Banner    map[PairIDs]banner.GetFilteredBannersResponse
	BannerFID map[int64][]banner.GetFilteredBannersResponse
	BannerTID map[int64][]banner.GetFilteredBannersResponse
	repo      banner.Repository
}

func NewCache(repo banner.Repository) cache.Cache {
	return &Cache{
		Banner:    make(map[PairIDs]banner.GetFilteredBannersResponse),
		BannerFID: make(map[int64][]banner.GetFilteredBannersResponse),
		BannerTID: make(map[int64][]banner.GetFilteredBannersResponse),
		repo:      repo,
	}
}

func (c *Cache) LoadCache() error {
	banners, err := c.repo.GetAllBanners()
	if err != nil {
		log.Println("Error in GetAllBanners:", err)
		return err
	}

	bannerMap := make(map[PairIDs]banner.GetFilteredBannersResponse)
	bannerMapFID := make(map[int64][]banner.GetFilteredBannersResponse)
	bannerMapTID := make(map[int64][]banner.GetFilteredBannersResponse)
	for _, bannerData := range *banners {
		tagsStr := strings.Split(strings.Trim(bannerData.TagIDs, "{}"), ",")
		for _, s := range tagsStr {
			tag, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Println("Error in LoadCache", err)
				continue
			}

			bannerMap[PairIDs{bannerData.FeatureID, tag}] = bannerData

			if _, ok := bannerMapFID[bannerData.FeatureID]; !ok {
				bannerMapFID[bannerData.FeatureID] = make([]banner.GetFilteredBannersResponse, 0)
			}
			bannerMapFID[bannerData.FeatureID] = append(bannerMapFID[bannerData.FeatureID], bannerData)

			if _, ok := bannerMapTID[tag]; !ok {
				bannerMapTID[tag] = make([]banner.GetFilteredBannersResponse, 0)
			}
			bannerMapTID[tag] = append(bannerMapTID[tag], bannerData)
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func(mutex *sync.RWMutex, bannerMap map[PairIDs]banner.GetFilteredBannersResponse) {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.Banner = bannerMap
	}(&c.mutex, bannerMap)

	go func(mutex *sync.RWMutex, bannerMapFID map[int64][]banner.GetFilteredBannersResponse) {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.BannerFID = bannerMapFID
	}(&c.mutexFID, bannerMapFID)

	go func(mutex *sync.RWMutex, bannerMapTID map[int64][]banner.GetFilteredBannersResponse) {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		c.BannerTID = bannerMapTID
	}(&c.mutexTID, bannerMapTID)

	wg.Wait()
	return nil
}

func (c *Cache) GetBanner(fid int64, tid int64) (banner.GetFilteredBannersResponse, error) {
	c.mutex.RLock()
	bannerData, ok := c.Banner[PairIDs{fid, tid}]
	c.mutex.RUnlock()

	if !ok {
		log.Printf("no data: PairIDs{fid, tid} =  {%d, %d}", fid, tid)
		return banner.GetFilteredBannersResponse{}, fmt.Errorf("no data: PairIDs{fid, tid}")
	}
	return bannerData, nil
}

func (c *Cache) GetBannersByTID(tid int64) ([]banner.GetFilteredBannersResponse, error) {
	c.mutexTID.RLock()
	bannersData, ok := c.BannerTID[tid]
	c.mutexTID.RUnlock()

	if !ok {
		log.Printf("no data: tid = %d", tid)
		return nil, fmt.Errorf("no data: tid")
	}
	return bannersData, nil
}

func (c *Cache) GetBannersByFID(fid int64) ([]banner.GetFilteredBannersResponse, error) {
	c.mutexFID.RLock()
	bannersData, ok := c.BannerFID[fid]
	c.mutexFID.RUnlock()

	if !ok {
		log.Printf("no data: fid = %d", fid)
		return nil, fmt.Errorf("no data: fid")
	}
	return bannersData, nil
}
