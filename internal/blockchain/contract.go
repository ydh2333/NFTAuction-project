package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
)

// AuctionContract 拍卖合约实例
type AuctionContract struct {
	client     *ethclient.Client // 区块链客户端
	address    common.Address    // 合约地址
	abi        abi.ABI           // 合约ABI
	privateKey string            // 私钥
	chainID    *big.Int          // 链ID
}

// 替换为你的合约ABI（从Solidity编译后获取）
const auctionABI = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "target",
				"type": "address"
			}
		],
		"name": "AddressEmptyCode",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_seller",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_duration",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_startPrice",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "_nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_nftId",
				"type": "uint256"
			}
		],
		"name": "createAuction",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "implementation",
				"type": "address"
			}
		],
		"name": "ERC1967InvalidImplementation",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "ERC1967NonPayable",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "FailedCall",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "InvalidInitialization",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "NotInitializing",
		"type": "error"
	},
	{
		"inputs": [],
		"name": "UUPSUnauthorizedCallContext",
		"type": "error"
	},
	{
		"inputs": [
			{
				"internalType": "bytes32",
				"name": "slot",
				"type": "bytes32"
			}
		],
		"name": "UUPSUnsupportedProxiableUUID",
		"type": "error"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "auctionId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "seller",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "duration",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "startPrice",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "startTokenAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "startTime",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "nftContract",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "nftId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "optTime",
				"type": "uint256"
			}
		],
		"name": "CreateAuction",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_auctionId",
				"type": "uint256"
			}
		],
		"name": "endAuction",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "auctionId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "winner",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "tokenAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "optTime",
				"type": "uint256"
			}
		],
		"name": "EndAuction",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "initialize",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint64",
				"name": "version",
				"type": "uint64"
			}
		],
		"name": "Initialized",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_auctionId",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "_tokenAddress",
				"type": "address"
			}
		],
		"name": "placeBid",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "uint256",
				"name": "auctionId",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "bidder",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"indexed": false,
				"internalType": "address",
				"name": "tokenAddress",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "optTime",
				"type": "uint256"
			}
		],
		"name": "PlaceBid",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_tokenAddress",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_priceFeedAddress",
				"type": "address"
			}
		],
		"name": "setPriceETHFeed",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "implementation",
				"type": "address"
			}
		],
		"name": "Upgraded",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newImplementation",
				"type": "address"
			},
			{
				"internalType": "bytes",
				"name": "data",
				"type": "bytes"
			}
		],
		"name": "upgradeToAndCall",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "admin",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "auctions",
		"outputs": [
			{
				"internalType": "address",
				"name": "seller",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "duration",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "startTime",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "startPrice",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "startTokenAddress",
				"type": "address"
			},
			{
				"internalType": "bool",
				"name": "ended",
				"type": "bool"
			},
			{
				"internalType": "uint256",
				"name": "highestBid",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "highestBidder",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "nftContract",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "nftId",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "tokenAddress",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_amount",
				"type": "uint256"
			},
			{
				"internalType": "address",
				"name": "_tokenAddress",
				"type": "address"
			}
		],
		"name": "calculateValue",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_tokenAddress",
				"type": "address"
			}
		],
		"name": "getChainlinkDataFeedLatestAnswer",
		"outputs": [
			{
				"internalType": "int256",
				"name": "",
				"type": "int256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "nextAuctionId",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "",
				"type": "bytes"
			}
		],
		"name": "onERC721Received",
		"outputs": [
			{
				"internalType": "bytes4",
				"name": "",
				"type": "bytes4"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "priceFeedDecimals",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "priceFeeds",
		"outputs": [
			{
				"internalType": "contract AggregatorV3Interface",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "proxiableUUID",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "UPGRADE_INTERFACE_VERSION",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// NewAuctionContract 初始化合约实例
func NewAuctionContract(cfg *config.BlockchainConfig) (*AuctionContract, error) {
	// 连接区块链节点
	client, err := ethclient.Dial(cfg.RPCEndpoint)
	if err != nil {
		log.Error().Err(err).Msg("区块链节点连接失败")
		return nil, err
	}

	// 解析合约ABI
	contractABI, err := abi.JSON(bytes.NewReader([]byte(auctionABI)))
	if err != nil {
		log.Error().Err(err).Msg("合约ABI解析失败")
		return nil, err
	}

	// 获取链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("获取链ID失败")
		return nil, err
	}

	return &AuctionContract{
		client:     client,
		address:    common.HexToAddress(cfg.ContractAddress),
		abi:        contractABI,
		privateKey: cfg.PrivateKey,
		chainID:    chainID,
	}, nil
}

// CreateAuction 调用合约创建拍卖（后端代创建）
func (c *AuctionContract) CreateAuction(seller common.Address, duration *big.Int, startPrice *big.Int, nftContract common.Address, nftId *big.Int) (string, error) {
	// 加载私钥
	privateKey, err := crypto.HexToECDSA(c.privateKey)
	if err != nil {
		log.Error().Err(err).Msg("私钥解析失败")
		return "", err
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error().Msg("error casting public key to ECDSA")
		return "", err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取nonce
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Error().Err(err).Msg("获取nonce失败")
		return "", err
	}

	// 获取燃气价格
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("获取燃气价格失败")
		return "", err
	}

	// 调用合约方法
	data, err := c.abi.Pack("createAuction", seller, duration, startPrice, nftContract, nftId)
	if err != nil {
		log.Error().Err(err).Msg("合约方法打包失败")
		return "", err
	}

	// 估算燃气限制
	gasLimit, err := c.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &c.address,
		Data: data,
	})
	if err != nil {
		log.Error().Err(err).Msg("估算燃气限制失败")
		return "", err
	}

	// 创建未签名交易（数据字段为 ABI 编码后的函数调用数据）
	tx := types.NewTransaction(
		nonce,
		c.address,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		data,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), privateKey)
	if err != nil {
		log.Error().Err(err).Msg("交易签名失败")
		return "", err
	}

	// 发送交易到区块链
	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Error().Err(err).Msg("交易发送失败")
		return "", err
	}
	log.Info().Msg("交易hash地址：" + string(signedTx.Hash().Hex()))

	// 等待交易回执
	receipt, err := waitForReceipt(c.client, signedTx.Hash())
	if err != nil {
		log.Error().Err(err).Msg("等待交易回执失败")
		return "", err
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("交易成功！区块号：%d，消耗Gas：%d\n", receipt.BlockNumber, receipt.GasUsed)
	} else {
		fmt.Println("交易失败！")
	}
	return string(signedTx.Hash().Hex()), nil
}

// PlaceBid 调用合约提交竞拍
func (c *AuctionContract) PlaceBid(auctionID, bidAmount *big.Int, tokenAddress common.Address, bidderPrivateKey string) (string, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(bidderPrivateKey)
	if err != nil {
		return "", err
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error().Msg("error casting public key to ECDSA")
		return "", err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取燃气价格
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("获取燃气价格失败")
		return "", err
	}

	// 获取nonce
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Error().Err(err).Msg("获取nonce失败")
		return "", err
	}

	// 调用合约方法
	data, err := c.abi.Pack("placeBid", auctionID, bidAmount, tokenAddress)
	if err != nil {
		log.Error().Err(err).Msg("合约方法打包失败")
		return "", err
	}

	// 估算燃气限制（智能合约调用需估算，而非固定 21000）
	gasLimit, err := c.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddress,
		To:    &c.address,
		Data:  data,
		Value: bidAmount,
	})
	if err != nil {
		log.Error().Err(err).Msg("估算燃气限制失败")
		return "", err
	}

	// 创建未签名交易, 判断是否用eth支付
	var tx *types.Transaction
	ethZeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	if tokenAddress == ethZeroAddress {
		tx = types.NewTransaction(
			nonce,
			c.address,
			bidAmount,
			gasLimit,
			gasPrice,
			data,
		)
	} else {
		tx = types.NewTransaction(
			nonce,
			c.address,
			big.NewInt(0),
			gasLimit,
			gasPrice,
			data,
		)
	}

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.chainID), privateKey)
	if err != nil {
		log.Error().Err(err).Msg("交易签名失败")
		return "", err
	}

	err = c.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Error().Err(err).Msg("交易发送失败")
		return "", err
	}
	log.Info().Msg("交易hash地址：" + string(signedTx.Hash().Hex()))

	// 等待交易回执
	receipt, err := waitForReceipt(c.client, signedTx.Hash())
	if err != nil {
		log.Error().Err(err).Msg("等待交易回执失败")
		return "", err
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("交易成功！区块号：%d，消耗Gas：%d\n", receipt.BlockNumber, receipt.GasUsed)
	} else {
		fmt.Println("交易失败！")
	}

	return string(signedTx.Hash().Hex()), nil
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}
