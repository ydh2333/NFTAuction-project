package models

import (
	"gorm.io/gorm"
)

// Bid 竞拍记录
type Bid struct {
	gorm.Model
	AuctionID     uint    `gorm:"not null;index" json:"auction_id"` // 关联拍卖ID
	Auction       Auction `gorm:"foreignKey:AuctionID" json:"auction,omitempty"`
	BidderAddress string  `gorm:"not null" json:"bidder_address"`  // 竞拍者钱包地址
	Amount        uint64  `gorm:"not null" json:"amount"`          // 竞拍金额
	TokenAddress  string  `gorm:"not null" json:"token_address"`   // 竞拍的代币类型
	IsWinning     bool    `gorm:"default:false" json:"is_winning"` // 是否为当前最高价
}
