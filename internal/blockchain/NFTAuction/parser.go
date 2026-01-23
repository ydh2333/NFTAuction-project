package NFTAuction

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/utils/logger"
)

// 处理AuctionCreated事件（拍卖创建）
func (l *Listener) handleAuctionCreated(log types.Log) error {
	var event struct {
		AuctionId         *big.Int
		Seller            common.Address
		Duration          *big.Int
		StartPrice        *big.Int
		StartTokenAddress common.Address
		StartTime         *big.Int
		NftContract       common.Address
		NftId             *big.Int
		OptTime           *big.Int
	}

	// 解析事件数据
	if err := l.abi.UnpackIntoInterface(&event, "CreateAuction", log.Data); err != nil {
		return logger.WrapError(err, "解析CreateAuction事件失败")
	}

	// 解析索引字段（topics）
	event.AuctionId = new(big.Int).SetBytes(log.Topics[1].Bytes())
	EndTime := new(big.Int).Add(event.StartTime, event.Duration)

	// 存入数据库
	auction := &models.Auction{
		ID:                event.AuctionId.Uint64(),
		CreatorAddress:    event.Seller.Hex(),
		Duration:          time.Duration(event.Duration.Uint64()) * time.Second,
		StartPrice:        event.StartPrice.Uint64(),
		StartTokenAddress: event.StartTokenAddress.Hex(),
		StartTime:         time.Unix(int64(event.StartTime.Uint64()), 0),
		EndTime:           time.Unix(int64(EndTime.Uint64()), 0),
		NFTContract:       event.NftContract.Hex(),
		NFTTokenID:        uint(event.NftId.Uint64()),
		OptTime:           time.Unix(int64(event.OptTime.Uint64()), 0),
	}

	// 保存拍卖数据
	auctionRepository := repository.NewAuctionRepository()
	if err := auctionRepository.Create(auction); err != nil {
		return logger.WrapError(err, "保存拍卖数据失败")
	}

	logger.Log.Info().Uint64("auction_id", auction.ID).Msg("同步拍卖创建事件成功")
	return nil
}

// 处理BidPlaced事件（出价）
func (l *Listener) handleBidPlaced(log types.Log) error {
	var event struct {
		AuctionId    *big.Int
		Bidder       common.Address
		Amount       *big.Int
		TokenAddress common.Address
		OptTime      *big.Int
	}

	if err := l.abi.UnpackIntoInterface(&event, "PlaceBid", log.Data); err != nil {
		return logger.WrapError(err, "解析PlaceBid事件失败")
	}

	// 解析索引字段
	event.AuctionId = new(big.Int).SetBytes(log.Topics[1].Bytes())

	// 1. 保存出价记录
	bid := &models.Bid{
		AuctionID:     event.AuctionId.Uint64(),
		BidderAddress: event.Bidder.Hex(),
		Amount:        event.Amount.Uint64(),
		TokenAddress:  event.TokenAddress.Hex(),
		OptTime:       time.Unix(int64(event.OptTime.Uint64()), 0),
	}

	// 2. 创建拍卖表记录，更新拍卖的当前最高价和出价者，二者需要保持数据一致性
	tx := repository.DB.Begin()
	if tx.Error != nil {
		return logger.WrapError(tx.Error, "开启事务失败")
	}

	bidRepository := repository.NewBidRepositoryWithTx(tx)
	if err := bidRepository.Create(bid); err != nil {
		return logger.WrapError(err, "保存出价记录失败")
	}
	auctionRepository := repository.NewAuctionRepositoryWithTx(tx)
	if err := auctionRepository.UpdateCurrentPrice(uint(bid.AuctionID), bid.Amount, bid.BidderAddress, bid.TokenAddress); err != nil {
		return logger.WrapError(err, "更新拍卖的当前最高价和出价者失败")
	}

	// 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		return logger.WrapError(err, "提交事务失败")
	}

	logger.Log.Info().Uint64("auction_id", bid.AuctionID).Str("bidder", bid.BidderAddress).Msg("同步出价事件成功")
	return nil
}

// 处理AuctionEnded事件（拍卖结束）
func (l *Listener) handleAuctionEnded(log types.Log) error {
	var event struct {
		AuctionId    *big.Int
		Winner       common.Address
		Amount       *big.Int
		TokenAddress common.Address
		OptTime      *big.Int
	}

	if err := l.abi.UnpackIntoInterface(&event, "AuctionEnded", log.Data); err != nil {
		return logger.WrapError(err, "解析AuctionEnded事件失败")
	}

	event.AuctionId = new(big.Int).SetBytes(log.Topics[1].Bytes())
	auctionID := uint(event.AuctionId.Uint64())

	tx := repository.DB.Begin()
	if tx.Error != nil {
		return logger.WrapError(tx.Error, "开启事务失败")
	}
	// 更新拍卖状态
	auctionRepository := repository.NewAuctionRepositoryWithTx(tx)
	auctionStatus := models.AuctionStatusEnded
	if err := auctionRepository.UpdateStatus(auctionID, auctionStatus); err != nil {
		return logger.WrapError(err, "更新拍卖状态失败")
	}

	// 更新拍卖的当前最高价和出价者
	bidRepository := repository.NewBidRepositoryWithTx(tx)
	if err := bidRepository.MarkWinningBid(auctionID); err != nil {
		return logger.WrapError(err, "标记获胜竞拍失败")
	}

	// 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		return logger.WrapError(err, "提交事务失败")
	}

	logger.Log.Info().Uint64("auction_id", uint64(auctionID)).Str("status", string(auctionStatus)).Msg("同步拍卖结束事件成功")
	return nil
}

// 定期检查进行中、过期拍卖
func (l *Listener) checkAuctionExpiry(ctx context.Context) error {

	ticker := time.NewTicker(time.Duration(l.pollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// 获取当前时间戳（秒）
			currentTime := time.Now().Unix()

			// 更新：开始时间 < 当前时间 → 进行中
			result := repository.DB.Model(&models.Auction{}).
				Where("status = ? AND start_time < FROM_UNIXTIME(?)", "pending", currentTime).
				Update("status", "active")

			if result.Error != nil {
				log.Error().Err(result.Error).Msg("检查过期拍卖失败")
			}

			// 更新：进行中且结束时间 < 当前时间 → 流拍
			result = repository.DB.Model(&models.Auction{}).
				Where("status = ? AND end_time < FROM_UNIXTIME(?)", "active", currentTime).
				Update("status", "ended")

			if result.Error != nil {
				logger.Log.Error().Err(result.Error).Msg("检查过期拍卖失败")
			} else if result.RowsAffected > 0 {
				logger.Log.Info().Int64("count", result.RowsAffected).Msg("更新过期流拍拍卖")
			}
		}
	}
}
