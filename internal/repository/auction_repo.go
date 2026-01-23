package repository

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/utils"
	"gorm.io/gorm"
)

type AuctionRepository interface {
	Create(auction *models.Auction) error
	GetByID(id uint) (*models.Auction, error)
	GetActiveAuctions() ([]*models.Auction, error)
	UpdateStatus(id uint, status models.AuctionStatus) error
	UpdateCurrentPrice(auctionID uint, HighestBid uint64, HighestBidder string, TokenAddress string) error
	GetAuctionCount() (int64, error)
	SearchAuctions(params AuctionSearchParams, sortParams SortParams, pageParams utils.PageParams) ([]AuctionDetail, error)
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

// getAuctionCount 获取拍卖总数量
func (r *auctionRepository) GetAuctionCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Auction{}).Count(&count).Error; err != nil {
		log.Error().Err(err).Msg("获取拍卖总数量失败")

		return 0, err
	}
	return count, nil
}

// 动态搜索参数
type AuctionSearchParams struct {
	Name          string
	TokenID       uint
	EndTimeMin    time.Time
	EndTimeMax    time.Time
	HighestBidMin uint64
	HighestBidMax uint64
	StartPriceMin uint64
	StartPriceMax uint64
	Status        models.AuctionStatus
}

// 封装搜索范围
func SearchAuctions(params AuctionSearchParams) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if params.Name != "" {
			tx = tx.Where("name LIKE ?", "%"+params.Name+"%")
		}
		if params.TokenID != 0 {
			tx = tx.Where("nft_token_id = ?", params.TokenID)
		}
		if !params.EndTimeMin.IsZero() {
			tx = tx.Where("end_time >= ?", params.EndTimeMin)
		}
		if !params.EndTimeMax.IsZero() {
			tx = tx.Where("end_time <= ?", params.EndTimeMax)
		}
		if params.HighestBidMin != 0 {
			tx = tx.Where("highest_bid >= ?", params.HighestBidMin)
		}
		if params.HighestBidMax != 0 {
			tx = tx.Where("highest_bid <= ?", params.HighestBidMax)
		}
		if params.StartPriceMin != 0 {
			tx = tx.Where("start_price >= ?", params.StartPriceMin)
		}
		return tx
	}
}

// 排序参数
type SortParams struct {
	Field string // 排序字段
	Dir   string // 排序方向
}

func SortAuctions(params SortParams) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// 允许的排序字段
		allowedFields := map[string]bool{
			"token_id":    true,
			"end_time":    true,
			"highest_bid": true,
			"start_price": true,
		}

		// 默认排序，按照起始价格降序
		if params.Field == "" && allowedFields[params.Field] {
			params.Field = "start_price"
			params.Dir = "desc"
		}

		if params.Dir != "asc" && params.Dir != "desc" {
			params.Dir = "desc"
		}

		return tx.Order(params.Field + " " + params.Dir)
	}
}

type AuctionDetail struct {
	ImageURL   string
	Name       string
	TokenID    string
	EndTime    time.Time
	HighestBid uint64
	StartPrice float64
	Status     models.AuctionStatus
}

func (r *auctionRepository) SearchAuctions(params AuctionSearchParams, sortParams SortParams, pageParams utils.PageParams) ([]AuctionDetail, error) {
	var AuctionDetails []AuctionDetail

	err := r.db.Table("auctions").
		Joins("JOIN nfts ON nfts.token_id = auctions.nft_token_id").
		Scopes(SearchAuctions(params), SortAuctions(sortParams), utils.Paginate(pageParams)).
		Select("nfts.image_url, nfts.name, nfts.token_id, auctions.end_time, auctions.highest_bid, auctions.start_price, auctions.status").
		Scan(&AuctionDetails).Error

	if err != nil {
		log.Error().Err(err).Msg("搜索拍卖失败")
		return nil, err
	}

	return AuctionDetails, nil
}
