package service

import (
	"github.com/rs/zerolog/log"

	"github.com/ydh2333/NFTAuction-project/internal/redis"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/utils"
)

type HomePageService interface {
	PlatformStatistics() (int, int)
	SearchAuctionsList(
		params repository.AuctionSearchParams,
		sortParams repository.SortParams,
		pageParams utils.PageParams,
	) ([]repository.AuctionDetail, error)
	GetTop5HotAuctions() ([]repository.AuctionDetail, error)
}

type homePageService struct {
	auctionRepo repository.AuctionRepository
	bidRepo     repository.BidRepository
}

func NewHomePageService() HomePageService {
	return &homePageService{
		auctionRepo: repository.NewAuctionRepository(),
		bidRepo:     repository.NewBidRepository(),
	}
}

func (h *homePageService) PlatformStatistics() (int, int) {
	auctionCount, err := h.auctionRepo.GetAuctionCount()
	log.Info().Int64("auctionCount", auctionCount).Msg("获取拍卖总数量")
	if err != nil {
		log.Error().Err(err).Msg("获取拍卖总数量失败")
		return 0, 0
	}
	bidCount, err := h.bidRepo.GetBidCount()
	if err != nil {
		return 0, 0
	}
	return int(auctionCount), int(bidCount)
}

func (h *homePageService) SearchAuctionsList(
	params repository.AuctionSearchParams,
	sortParams repository.SortParams,
	pageParams utils.PageParams,
) ([]repository.AuctionDetail, error) {
	return h.auctionRepo.SearchAuctions(params, sortParams, pageParams)
}

func (h *homePageService) GetTop5HotAuctions() ([]repository.AuctionDetail, error) {
	// 1. 从redis中获取热门拍卖ID
	auctionIDs, err := redis.GetTop5HotAuctions()
	if err != nil {
		log.Error().Err(err).Msg("获取热门拍卖ID失败")
		return nil, err
	}

	// 2. 根据热门拍卖ID获取拍卖详情
	auctionRepo := repository.NewAuctionRepository()
	auctionDetails, err := auctionRepo.GetAuctionsByIDs(auctionIDs)
	if err != nil {
		log.Error().Err(err).Msg("获取热门拍卖失败")
		return nil, err
	}

	return auctionDetails, nil
}
