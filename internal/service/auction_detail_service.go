package service

import (
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
)

type AuctionDetailService interface {
	GetAuctionDetail(auctionId uint) ([]*models.Bid, error)
}

type auctionDetailService struct {
	bidRepo repository.BidRepository
}

func NewAuctionDetailService() AuctionDetailService {
	return &auctionDetailService{
		bidRepo: repository.NewBidRepository(),
	}
}

func (a *auctionDetailService) GetAuctionDetail(auctionId uint) ([]*models.Bid, error) {

	return a.bidRepo.GetByAuctionID(auctionId)
}
