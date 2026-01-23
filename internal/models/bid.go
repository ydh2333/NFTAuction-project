package models

import (
	"time"
)

// Bid 竞拍记录
type Bid struct {
	ID      uint `gorm:"primarykey"`
	OptTime time.Time

	BidderAddress string  `gorm:"not null" json:"bidder_address"`   // 竞拍者钱包地址
	Amount        uint64  `gorm:"not null" json:"amount"`           // 竞拍金额
	TokenAddress  string  `gorm:"not null" json:"token_address"`    // 竞拍的代币类型
	IsWinning     bool    `gorm:"default:false" json:"is_winning"`  // 是否为最高价，默认为false，结束后标记为true
	AuctionID     uint64  `gorm:"not null;index" json:"auction_id"` // 关联拍卖ID
	Auction       Auction `gorm:"foreignKey:AuctionID" json:"auction"`
}
