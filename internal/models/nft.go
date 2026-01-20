package models

import (
	"time"

	"gorm.io/gorm"
)

// NFT NFT基本信息
type NFT struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	TokenID         string         `gorm:"uniqueIndex;not null" json:"token_id"` // 区块链上的NFT TokenID
	ContractAddress string         `gorm:"not null" json:"contract_address"`     // NFT合约地址
	OwnerAddress    string         `gorm:"not null" json:"owner_address"`        // 当前所有者钱包地址
	Name            string         `gorm:"not null" json:"name"`                 // NFT名称
	Description     string         `json:"description"`                          // NFT描述
	ImageURL        string         `json:"image_url"`                            // NFT图片链接
}
