package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/utils"
)

func main() {
	// 初始化日志
	// logger.InitLogger()
	cfg := config.LoadConfig()

	repository.InitDB(&cfg.MySQL)

	// nft := &models.NFT{
	// 	TokenID:         1,
	// 	ContractAddress: "0x123",
	// 	OwnerAddress:    "0x456",
	// 	Name:            "NFT",
	// 	ImageURL:        "https://example.com/nft.jpg",
	// }

	// nftRepo := repository.NewNFTRepository()
	// nftRepo.Create(nft)

	// auction := &models.Auction{
	// 	CreatorAddress:    "0x123",
	// 	Duration:          3600,
	// 	StartTime:         time.Now(),
	// 	EndTime:           time.Now().Add(time.Hour),
	// 	StartPrice:        100,
	// 	StartTokenAddress: "0x456",
	// 	Status:            models.AuctionStatusPending,
	// 	HighestBidder:     "",
	// 	HighestBid:        0,
	// 	TokenAddress:      "0x789",
	// 	NFTTokenID:        1,
	// }

	auctionRepo := repository.NewAuctionRepository()
	// auctionRepo.Create(auction)

	// nftRepository := repository.NewNFTRepository()
	// nftDetails, err := nftRepository.GetNFTByOwnerAddress("0x51ccc58AE0a621b78196CcE2e01920dd6E5be38b")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Info().Interface("nftDetails", nftDetails).Msg("查询NFT成功")

	// nft, err := nftRepository.GetNFTByTokenID(21)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Info().Interface("nft", nft).Msg("查询NFT成功")
	params := repository.AuctionSearchParams{
		Name: "My first NFT",
		// TokenID: 21,
		// EndTimeMin:    time.Now(),
		// EndTimeMax:    time.Now(),
		HighestBidMin: 100,
		HighestBidMax: 200,
		StartPriceMin: 100,
		StartPriceMax: 200,
		Status:        models.AuctionStatusPending,
	}
	sortParams := repository.SortParams{
		Field: "highest_bid",
		Dir:   "desc",
	}
	pageParams := utils.PageParams{
		Page: 1,
		Size: 10,
	}
	auctions, err := auctionRepo.SearchAuctions(params, sortParams, pageParams)
	if err != nil {
		panic(err)
	}
	log.Info().Interface("auction", auctions).Msg("查询拍卖成功")

}
