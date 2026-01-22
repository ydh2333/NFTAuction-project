package repository

import (
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"gorm.io/gorm"
)

type BidRepository interface {
	Create(bid *models.Bid) error
	GetByAuctionID(auctionID uint) ([]*models.Bid, error)
	GetHighestBidByAuctionID(auctionID uint) (*models.Bid, error)
	MarkWinningBid(bidID uint) error
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository() BidRepository {
	return &bidRepository{db: DB}
}

func NewBidRepositoryWithTx(tx *gorm.DB) BidRepository {
	return &bidRepository{db: tx}
}

// Create 创建竞拍记录
func (r *bidRepository) Create(bid *models.Bid) error {
	if err := r.db.Create(bid).Error; err != nil {
		log.Error().Err(err).Msg("创建竞拍记录失败")
		return err
	}
	return nil
}

// GetByAuctionID 根据拍卖ID查询竞拍记录
func (r *bidRepository) GetByAuctionID(auctionID uint) ([]*models.Bid, error) {
	var bids []*models.Bid
	if err := r.db.Where("auction_id = ?", auctionID).Find(&bids).Error; err != nil {
		log.Error().Err(err).Msg("查询竞拍记录失败")
		return nil, err
	}
	return bids, nil
}

// GetHighestBidByAuctionID 根据拍卖ID查询最高竞拍
func (r *bidRepository) GetHighestBidByAuctionID(auctionID uint) (*models.Bid, error) {
	var bid models.Bid
	if err := r.db.Where("auction_id = ?", auctionID).Order("amount DESC").First(&bid).Error; err != nil {
		log.Error().Err(err).Msg("查询最高竞拍失败")
		return nil, err
	}
	return &bid, nil
}

// MarkWinningBid 标记获胜竞拍，拍卖结束时更新
func (r *bidRepository) MarkWinningBid(bidID uint) error {
	if err := r.db.Model(&models.Bid{}).Where("id = ?", bidID).Update("is_winning", true).Error; err != nil {
		log.Error().Err(err).Msg("标记获胜竞拍失败")
		return err
	}
	return nil
}
