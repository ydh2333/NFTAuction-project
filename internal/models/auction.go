package models

import (
	"time"
)

// AuctionStatus 拍卖状态
type AuctionStatus string

const (
	AuctionStatusPending   AuctionStatus = "pending"   // 未开始
	AuctionStatusActive    AuctionStatus = "active"    // 进行中
	AuctionStatusEnded     AuctionStatus = "ended"     // 已结束
	AuctionStatusCancelled AuctionStatus = "cancelled" // 已取消
)

// Auction 拍卖信息
type Auction struct {
	ID      uint64 `gorm:"primarykey" json:"id"`
	OptTime time.Time

	CreatorAddress    string        `gorm:"not null" json:"creator_address"`          // 拍卖创建者钱包地址
	Duration          time.Duration `gorm:"not null" json:"duration"`                 // 拍卖持续时间
	StartTime         time.Time     `gorm:"not null" json:"start_time"`               // 拍卖开始时间
	EndTime           time.Time     `gorm:"not null" json:"end_time"`                 // 拍卖结束时间
	StartPrice        uint64        `gorm:"not null" json:"start_price"`              // 起拍价（单位：ETH/USDT等）
	StartTokenAddress string        `gorm:"not null" json:"start_token_address"`      // 起始货币类型
	Status            AuctionStatus `gorm:"not null;default:'pending'" json:"status"` // 拍卖状态
	HighestBidder     *uint         `gorm:"not null;default:0" json:"highest_bidder"` // 当前最高价出价者
	HighestBid        uint64        `gorm:"not null;default:0" json:"highest_bid"`    // 当前最高价
	TokenAddress      string        `gorm:"not null" json:"token_address"`            // 拍卖货币类型
	NFTID             uint          `gorm:"not null;index" json:"nft_id"`             // 关联NFT的ID
	NFTContract       string        `gorm:"not null" json:"nft_contract"`             // NFT合约地址

}
