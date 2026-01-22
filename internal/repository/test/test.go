package main

import (
	"time"

	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
)

func main() {
	// 初始化日志
	// logger.InitLogger()
	cfg := config.LoadConfig()

	repository.InitDB(&cfg.MySQL)

	nft := &models.NFT{
		TokenID:         1,
		ContractAddress: "0x123",
		OwnerAddress:    "0x456",
		Name:            "NFT",
		ImageURL:        "https://example.com/nft.jpg",
	}

	nftRepo := repository.NewNFTRepository()
	nftRepo.Create(nft)

	auction := &models.Auction{
		CreatorAddress:    "0x123",
		Duration:          3600,
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		StartPrice:        100,
		StartTokenAddress: "0x456",
		Status:            models.AuctionStatusPending,
		HighestBidder:     "",
		HighestBid:        0,
		TokenAddress:      "0x789",
		NFTID:             1,
	}

	auctionRepo := repository.NewAuctionRepository()
	auctionRepo.Create(auction)

}
