package repository

import (
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"gorm.io/gorm"
)

type NFTRepository interface {
	Create(nft *models.NFT) error
	GetNFTByTokenID(id uint) (*models.NFT, error)
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
func (r *nftRepository) GetNFTByTokenID(id uint) (*models.NFT, error) {
	var nft models.NFT

	if err := r.db.First(&nft, id).Error; err != nil {
		log.Error().Err(err).Uint("nft_id", id).Msg("查询NFT失败")
		return nil, err
	}

	return &nft, nil
}
