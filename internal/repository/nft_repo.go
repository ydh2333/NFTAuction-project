package repository

import (
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"gorm.io/gorm"
)

type NFTRepository interface {
	Create(nft *models.NFT) error
	GetNFTByTokenID(tokenID uint) (*models.NFT, error)
	GetNFTByOwnerAddress(OwnerAddress string) ([]NftDetail, error)
}

type nftRepository struct {
	db *gorm.DB
}

func NewNFTRepository() NFTRepository {
	return &nftRepository{db: DB}
}

func (r *nftRepository) Create(nft *models.NFT) error {
	if err := r.db.Create(nft).Error; err != nil {
		log.Error().Err(err).Msg("创建NFT失败")
		return err
	}
	return nil
}

// GetNFTByTokenID 根据TokenID查询NFT
func (r *nftRepository) GetNFTByTokenID(tokenID uint) (*models.NFT, error) {
	var nft models.NFT

	if err := r.db.First(&nft).Where("token_id=?", tokenID).Error; err != nil {
		log.Error().Err(err).Uint("nft_id", tokenID).Msg("查询NFT失败")
		return nil, err
	}

	return &nft, nil
}

type NftDetail struct {
	ImageURL   string
	Name       string
	TokenID    string
	StartPrice float64
	Status     models.AuctionStatus
}

// GetNFTByOwnerAddress 查询个人NFT拍卖列表
func (r *nftRepository) GetNFTByOwnerAddress(OwnerAddress string) ([]NftDetail, error) {
	var nftDetails []NftDetail
	err := r.db.Table("nfts").
		Joins("LEFT JOIN auctions ON nfts.token_id = auctions.nft_token_id").
		Where("nfts.owner_address = ?", OwnerAddress).
		Select("nfts.image_url, nfts.name, nfts.token_id, auctions.start_price, auctions.status").
		Scan(&nftDetails).Error

	if err != nil {
		log.Error().Err(err).Str("owner_address", OwnerAddress).Msg("查询NFT失败")
		return nil, err
	}
	return nftDetails, nil
}
