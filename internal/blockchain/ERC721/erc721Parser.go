package ERC721

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
)

// getTokenURI 调用合约tokenURI方法获取元数据链接
func (l *ERC721Listener) getTokenURI(ctx context.Context, tokenId uint64) (string, error) {
	tokenIdBig := big.NewInt(int64(tokenId))
	// 打包调用参数
	data, err := l.abi.Pack("tokenURI", tokenIdBig)
	if err != nil {
		return "", fmt.Errorf("打包tokenURI参数失败: %w", err)
	}

	// 构造合约调用消息
	callMsg := ethereum.CallMsg{
		To:   &l.contractAddr,
		Data: data,
	}

	// 执行合约调用
	result, err := l.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return "", fmt.Errorf("调用tokenURI失败: %w", err)
	}

	// 解包返回值（string类型）
	var tokenURI string
	if err := l.abi.UnpackIntoInterface(&tokenURI, "tokenURI", result); err != nil {
		return "", fmt.Errorf("解包tokenURI结果失败: %w", err)
	}

	return tokenURI, nil
}

// resolveMetadata 解析元数据链接（支持IPFS和HTTP）
func (l *ERC721Listener) resolveMetadata(ctx context.Context, tokenURI string) (*NFTMetadata, error) {
	// 处理IPFS链接（转换为HTTP可访问的链接）
	resolvedURI := tokenURI
	if strings.HasPrefix(tokenURI, "ipfs://") {
		resolvedURI = "https://ipfs.io/ipfs/" + strings.TrimPrefix(tokenURI, "ipfs://")
	}

	// 发送HTTP请求获取元数据
	req, err := http.NewRequestWithContext(ctx, "GET", resolvedURI, nil)
	if err != nil {
		return nil, fmt.Errorf("构建元数据请求失败: %w", err)
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求元数据失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取元数据响应失败: %w", err)
	}

	// 解析JSON元数据
	var metadata NFTMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("解析元数据JSON失败: %w, 原始数据: %s", err, string(body))
	}

	return &metadata, nil
}

// handleSafeMint 处理safeMint事件（解析Transfer日志）
func (l *ERC721Listener) handleSafeMint(ctx context.Context, logEntry types.Log) error {
	// 验证Transfer事件的Topics数量
	if len(logEntry.Topics) < 4 {
		return fmt.Errorf("日志Topics不足，无法解析Transfer事件")
	}

	// 解析tokenID（uint256转uint64，如需高精度可改用big.Int）
	var tokenIdHash common.Hash
	tokenIdHash.SetBytes(logEntry.Topics[3].Bytes())
	tokenIdBig := tokenIdHash.Big()
	tokenId := tokenIdBig.Uint64()

	// 解析接收钱包地址
	walletAddr := common.HexToAddress(logEntry.Topics[2].Hex()).Hex()

	blockTimeUnix, err := l.getBlockTime(ctx, logEntry.BlockNumber)
	if err != nil {
		return fmt.Errorf("获取区块时间失败: %w", err)
	}

	// 1. 获取tokenURI
	tokenURI, err := l.getTokenURI(ctx, tokenId)
	if err != nil {
		log.Warn().Err(err).Uint64("tokenId", tokenId).Msg("获取tokenURI失败，仅输出基础信息")
		// 输出基础信息（无元数据）
		l.printNFTInfo(&NFTInfo{
			TokenID:      tokenId,
			ContractAddr: l.contractAddr.Hex(),
			WalletAddr:   walletAddr,
			TxHash:       logEntry.TxHash.Hex(),
		})
		return nil
	}

	// 2. 解析元数据
	metadata, err := l.resolveMetadata(ctx, tokenURI)
	if err != nil {
		log.Warn().Err(err).Uint64("tokenId", tokenId).Str("tokenURI", tokenURI).Msg("解析元数据失败")
		// 输出基础信息（无元数据）
		l.printNFTInfo(&NFTInfo{
			TokenID:      tokenId,
			ContractAddr: l.contractAddr.Hex(),
			WalletAddr:   walletAddr,
			TxHash:       logEntry.TxHash.Hex(),
		})
		return nil
	}

	nft := &models.NFT{
		TokenID:         uint(tokenId),
		ContractAddress: l.contractAddr.Hex(),
		OwnerAddress:    walletAddr,
		Name:            metadata.Name,
		Description:     metadata.Description,
		ImageURL:        metadata.Image,
		OptTime:         time.Unix(int64(blockTimeUnix), 0),
	}

	nftRepository := repository.NewNFTRepository()
	if err := nftRepository.Create(nft); err != nil {
		return err
	}

	return nil
}

// getBlockTime 根据区块号获取区块时间（格式化字符串 + Unix时间戳）
func (l *ERC721Listener) getBlockTime(ctx context.Context, blockNumber uint64) (uint64, error) {
	header, err := l.client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, fmt.Errorf("获取区块头失败: %w", err)
	}

	blockTimeUnix := header.Time

	return blockTimeUnix, nil
}

// printNFTInfo 格式化输出NFT信息
func (l *ERC721Listener) printNFTInfo(info *NFTInfo) {
	// 输出基础信息
	log.Info().
		Str("交易哈希", info.TxHash).
		Str("合约地址", info.ContractAddr).
		Str("接收钱包", info.WalletAddr).
		Uint64("TokenID", info.TokenID).
		Msg("捕获NFT safeMint事件")

	// 输出元数据（如果存在）
	if info.Metadata != nil {
		log.Info().
			Str("NFT名称", info.Metadata.Name).
			Str("NFT描述", info.Metadata.Description).
			Str("图片链接", info.Metadata.Image).
			Uint64("TokenID", info.TokenID).
			Msg("NFT元数据信息")
	}
}
