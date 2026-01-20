package models

import (
	"time"

	"gorm.io/gorm"
)

// Bid 竞拍记录
type Bid struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	AuctionID     uint           `gorm:"not null;index" json:"auction_id"` // 关联拍卖ID
	Auction       Auction        `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	BidderAddress string         `gorm:"not null" json:"bidder_address"`  // 竞拍者钱包地址
	Amount        uint           `gorm:"not null" json:"amount"`          // 竞拍金额
	IsWinning     bool           `gorm:"default:false" json:"is_winning"` // 是否为当前最高价
}
