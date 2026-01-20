package models

import (
	"time"

	"gorm.io/gorm"
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
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	NFTID          uint           `gorm:"not null;index" json:"nft_id"`             // 关联NFT的ID
	NFT            NFT            `gorm:"foreignKey:NFTID" json:"nft"`              // 关联NFT详情
	CreatorAddress string         `gorm:"not null" json:"creator_address"`          // 拍卖创建者钱包地址
	StartPrice     float64        `gorm:"not null" json:"start_price"`              // 起拍价（单位：ETH/USDT等）
	CurrentPrice   float64        `gorm:"not null;default:0" json:"current_price"`  // 当前最高价
	StartTime      time.Time      `gorm:"not null" json:"start_time"`               // 拍卖开始时间
	EndTime        time.Time      `gorm:"not null" json:"end_time"`                 // 拍卖结束时间
	Status         AuctionStatus  `gorm:"not null;default:'pending'" json:"status"` // 拍卖状态
	WinningBidID   *uint          `json:"winning_bid_id,omitempty"`                 // 获胜竞拍ID（结束后赋值）
}
