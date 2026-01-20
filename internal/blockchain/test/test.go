package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/blockchain"
	"github.com/ydh2333/NFTAuction-project/utils/logger"
)

func main() {
	logger.InitLogger()
	cfg := config.LoadConfig()

	blockchainClient, err := blockchain.NewAuctionContract(&cfg.Blockchain)
	if err != nil {
		log.Fatal().Err(err).Msg("创建合约实例失败")
		return
	}

	// tx, err := blockchainClient.CreateAuction(
	// 	common.HexToAddress("0x51ccc58ae0a621b78196cce2e01920dd6e5be38b"),
	// 	big.NewInt(100000),
	// 	big.NewInt(100),
	// 	common.HexToAddress("0x8174da3510e4C0373db82b92AB7949AfF75e7C25"),
	// 	big.NewInt(3),
	// )
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("创建拍卖失败")
	// 	return
	// }
	// log.Info().Str("txHash", tx).Msg("创建拍卖成功")

	tx, err := blockchainClient.PlaceBid(
		big.NewInt(1),
		big.NewInt(10000),
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		"d71a701e75b49c9a337ac20bacf15ccf62b92b86fee94cb6f5bc0240453f4f64",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("竞拍失败")
		return
	}
	log.Info().Str("txHash", tx).Msg("竞拍成功")
}
