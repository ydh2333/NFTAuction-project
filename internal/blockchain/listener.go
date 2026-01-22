package blockchain

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/utils/logger"
	"golang.org/x/sync/errgroup"
)

type Listener struct {
	client       *ethclient.Client
	abi          abi.ABI
	contractAddr common.Address
	startBlock   uint64
	pollInterval int64
}

// 初始化监听器
func NewListener(cfg *config.BlockchainConfig) (*Listener, error) {
	// 连接以太坊RPC
	client, err := ethclient.Dial(cfg.WSRpcEndpoint)
	if err != nil {
		return nil, err
	}

	// 解析ABI
	parsedABI, err := abi.JSON(strings.NewReader(AuctionContractABI))
	if err != nil {
		return nil, err
	}
	return &Listener{
		client:       client,
		abi:          parsedABI,
		contractAddr: common.HexToAddress(cfg.ContractAddr),
		startBlock:   cfg.StartBlock,
		pollInterval: int64(cfg.PollInterval),
	}, nil
}

// 启动监听（非阻塞，后台运行）
func (l *Listener) Start(ctx context.Context) error {
	// 1. 创建带上下文的errgroup，用于管理多个协程
	eg, ctx := errgroup.WithContext(ctx)

	// 2. 启动3个事件监听协程（AuctionCreated/BidPlaced/AuctionEnded）
	eg.Go(func() error {
		return l.listenEvent(ctx, "CreateAuction", l.handleAuctionCreated)
	})

	eg.Go(func() error {
		return l.listenEvent(ctx, "PlaceBid", l.handleBidPlaced)
	})

	eg.Go(func() error {
		return l.listenEvent(ctx, "EndAuction", l.handleAuctionEnded)
	})

	// 3. 启动拍卖过期检查协程（兜底逻辑）
	// 监听区块高度，更新拍卖状态（防止合约未触发AuctionEnded的情况）
	eg.Go(func() error {
		return l.checkAuctionExpiry(ctx)
	})

	// 4. 等待所有协程完成，返回第一个出错的错误
	return eg.Wait()
}

// 通用事件监听逻辑
func (l *Listener) listenEvent(ctx context.Context, eventName string, handler func(log types.Log) error) error {
	event := l.abi.Events[eventName]
	if event.ID == (common.Hash{}) {
		log.Error().Str("event", eventName).Msg("事件不存在")
		return logger.NewErrorf("事件%s不存在", eventName)
	}

	// 过滤条件：合约地址 + 事件签名
	query := ethereum.FilterQuery{
		Addresses: []common.Address{l.contractAddr},
		Topics:    [][]common.Hash{{event.ID}},
	}

	logs := make(chan types.Log)
	sub, err := l.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Error().Err(err).Str("event", eventName).Msg("订阅事件失败")
		return logger.WrapError(err, "订阅事件%s失败", eventName)
	}
	defer sub.Unsubscribe()

	// logger.Log.Info().Str("event", eventName).Msg("开始监听事件")
	log.Info().Str("event", eventName).Msg("开始监听事件")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			// logger.Log.Error().Err(err).Str("event", eventName).Msg("事件订阅出错，重试中...")
			log.Error().Err(err).Str("event", eventName).Msg("事件订阅出错，重试中...")
			// 重试订阅
			sub.Unsubscribe()
			sub, err = l.client.SubscribeFilterLogs(ctx, query, logs)
			if err != nil {
				log.Error().Err(err).Str("event", eventName).Msg("重试订阅事件失败")
				return logger.WrapError(err, "重试订阅事件%s失败", eventName)
			}
		case log1 := <-logs:
			// logger.Log.Debug().Str("event", eventName).Str("tx_hash", log.TxHash.Hex()).Msg("收到事件")
			log.Info().Str("event", eventName).Str("tx_hash", log1.TxHash.Hex()).Msg("收到事件")
			if err := handler(log1); err != nil {
				log.Error().Err(err).Str("event", eventName).Str("tx_hash", log1.TxHash.Hex()).Msg("处理事件失败")
				// logger.Log.Error().Err(err).Str("event", eventName).Str("tx_hash", log.TxHash.Hex()).Msg("处理事件失败")
			}
		}
	}
}
