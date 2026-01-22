package repository

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"gorm.io/gorm"
)

type AuctionRepository interface {
	Create(auction *models.Auction) error
	GetByID(id uint) (*models.Auction, error)
	GetActiveAuctions() ([]*models.Auction, error)
	UpdateStatus(id uint, status models.AuctionStatus) error
	UpdateCurrentPrice(auctionID uint, HighestBid uint64, HighestBidder string, TokenAddress string) error
}

// auctionRepository 实现AuctionRepository
type auctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository 创建拍卖仓库实例
func NewAuctionRepository() AuctionRepository {
	return &auctionRepository{db: DB}
}

func NewAuctionRepositoryWithTx(tx *gorm.DB) AuctionRepository {
	return &auctionRepository{db: tx}
}

// Create 创建拍卖记录
func (r *auctionRepository) Create(auction *models.Auction) error {
	if err := r.db.Create(auction).Error; err != nil {
		log.Error().Err(err).Msg("创建拍卖记录失败")
		return err
	}
	return nil
}

// GetByID 根据ID查询拍卖
func (r *auctionRepository) GetByID(id uint) (*models.Auction, error) {
	var auction models.Auction
	if err := r.db.Preload("NFT").First(&auction, id).Error; err != nil {
		log.Error().Err(err).Uint("auction_id", id).Msg("查询拍卖失败")
		return nil, err
	}
	return &auction, nil
}

// GetActiveAuctions 查询进行中的拍卖
func (r *auctionRepository) GetActiveAuctions() ([]*models.Auction, error) {
	var auctions []*models.Auction
	now := time.Now()
	if err := r.db.Preload("NFT").
		Where("status = ?", models.AuctionStatusActive).
		Where("start_time <= ?", now).
		Where("end_time >= ?", now).
		Find(&auctions).Error; err != nil {
		log.Error().Err(err).Msg("查询进行中拍卖失败")
		return nil, err
	}
	return auctions, nil
}

// UpdateStatus 更新拍卖状态
func (r *auctionRepository) UpdateStatus(id uint, status models.AuctionStatus) error {
	if err := r.db.Model(&models.Auction{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		log.Error().Err(err).Uint("auction_id", id).Msg("更新拍卖状态失败")
		return err
	}
	return nil
}

// UpdateCurrentPrice 更新拍卖当前最高价
func (r *auctionRepository) UpdateCurrentPrice(auctionID uint, HighestBid uint64, HighestBidder string, TokenAddress string) error {
	if err := r.db.Model(&models.Auction{}).
		Where("id = ?", auctionID).
		Updates(map[string]interface{}{
			"highest_bid":    HighestBid,
			"highest_bidder": HighestBidder,
			"token_address":  TokenAddress,
		}).Error; err != nil {
		log.Error().Err(err).Uint("auction_id", auctionID).Msg("更新拍卖当前价失败")
		return err
	}
	return nil
}
