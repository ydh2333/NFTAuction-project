package models

import (
	"time"
)

// NFT NFT基本信息
type NFT struct {
	ID      uint `gorm:"primarykey"`
	OptTime time.Time

	TokenID         uint   `gorm:"type:varchar(64);uniqueIndex;not null" json:"token_id"` // 区块链上的NFT TokenID
	ContractAddress string `gorm:"type:varchar(64);not null" json:"contract_address"`     // 合约地址，长度限制64字符哈希
	OwnerAddress    string `gorm:"type:varchar(64);not null" json:"owner_address"`        // 钱包地址
	Name            string `gorm:"type:varchar(255);not null" json:"name"`                // NFT名称
	Description     string `gorm:"type:text" json:"description"`                          // NFT描述
	ImageURL        string `gorm:"type:varchar(512)" json:"image_url"`                    // NFT图片链接
}
