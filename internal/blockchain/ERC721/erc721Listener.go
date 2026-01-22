package ERC721

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
)

// NFTMetadata NFT元数据结构体（适配主流ERC721元数据标准）
type NFTMetadata struct {
	Name        string `json:"name"`        // NFT名称
	Description string `json:"description"` // NFT描述
	Image       string `json:"image"`       // NFT图片链接
}

// NFTInfo 整合后的NFT完整信息
type NFTInfo struct {
	TokenID      uint64       // NFT代币ID
	ContractAddr string       // 合约地址
	WalletAddr   string       // 接收NFT的钱包地址
	Metadata     *NFTMetadata // NFT元数据
	TxHash       string       // 交易哈希
}

// ERC721Listener ERC721监听器
type ERC721Listener struct {
	client       *ethclient.Client // 以太坊RPC客户端
	abi          abi.ABI           // 解析后的ERC721 ABI
	contractAddr common.Address    // 监听的合约地址
	zeroAddr     common.Address    // 零地址（过滤safeMint）
	httpClient   *http.Client      // 解析元数据的HTTP客户端
}

// NewERC721Listener 初始化监听器
func NewERC721Listener(cfg *config.BlockchainConfig) (*ERC721Listener, error) {
	// 1. 连接以太坊RPC节点（如Infura、Alchemy或自建节点）
	client, err := ethclient.Dial(cfg.WSRpcEndpoint)
	if err != nil {
		return nil, err
	}

	// 2. 解析ERC721 ABI
	parsedABI, err := abi.JSON(strings.NewReader(ERC721ABI))
	if err != nil {
		return nil, err
	}

	// 3. 初始化HTTP客户端（设置超时，避免阻塞）
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &ERC721Listener{
		client:       client,
		abi:          parsedABI,
		contractAddr: common.HexToAddress(cfg.ERC721ContractAddr),
		zeroAddr:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		httpClient:   httpClient,
	}, nil
}

// StartListening 启动监听safeMint事件
func (l *ERC721Listener) StartListeningSafeMint(ctx context.Context) error {
	// 获取Transfer事件的ID（用于过滤日志）
	transferEvent, ok := l.abi.Events["Transfer"]
	if !ok {
		return fmt.Errorf("ABI中未找到Transfer事件")
	}
	eventID := transferEvent.ID

	// 构造日志过滤条件：
	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{l.contractAddr}, // 仅监听目标合约
		Topics: [][]common.Hash{
			{eventID},                            // Topics[0] = Transfer事件ID
			{common.HexToHash(l.zeroAddr.Hex())}, // Topics[1] = from地址（零地址，过滤safeMint）
			nil,                                  // Topics[2] = to地址（任意）
			nil,                                  // Topics[3] = tokenId（任意）
		},
	}

	// 订阅日志
	logs := make(chan types.Log)
	sub, err := l.client.SubscribeFilterLogs(ctx, filterQuery, logs)
	if err != nil {
		return fmt.Errorf("订阅日志失败: %w", err)
	}
	defer sub.Unsubscribe()

	log.Info().Str("监听合约", l.contractAddr.Hex()).Msg("开始监听ERC721 safeMint事件...")

	// 循环处理日志
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("监听停止：上下文已关闭")
			return ctx.Err()
		case err := <-sub.Err():
			// 订阅出错时重试
			log.Error().Err(err).Msg("订阅出错，2秒后重试")
			sub.Unsubscribe()
			time.Sleep(2 * time.Second)
			sub, err = l.client.SubscribeFilterLogs(ctx, filterQuery, logs)
			if err != nil {
				return fmt.Errorf("重试订阅失败: %w", err)
			}
		case logEntry := <-logs:
			// 处理单个日志条目
			if err := l.handleSafeMint(ctx, logEntry); err != nil {
				log.Error().Err(err).Str("交易哈希", logEntry.TxHash.Hex()).Msg("处理safeMint事件失败")
			}
		}
	}
}
