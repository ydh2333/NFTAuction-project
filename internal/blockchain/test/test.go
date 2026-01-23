package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/blockchain/ERC721"
	"github.com/ydh2333/NFTAuction-project/internal/blockchain/NFTAuction"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/utils/logger"
)

func main() {
	// 1. 初始化配置
	cfg := config.LoadConfig()

	// 2. 初始化日志
	logger.InitLogger()

	// 3. 初始化数据库
	repository.InitDB(&cfg.MySQL)

	// 创建全局上下文,统一管理所有监听器
	ctx, cancel := context.WithCancel(context.Background())
	// 确保程序退出时执行cancel，且只执行一次
	defer func() {
		cancel()
		log.Info().Msg("全局上下文已关闭，所有监听器停止")
	}()

	// 初始化erc721SafeMint监听器
	erc721Listener, err := ERC721.NewERC721Listener(&cfg.Blockchain)
	if err != nil {
		log.Fatal().Err(err).Msg("初始化监听器失败")
	}

	// 启动ERC721监听器（后台协程，避免阻塞主线程）
	go func() {
		log.Info().Msg("启动ERC721 safeMint监听器")
		if err := erc721Listener.StartListeningSafeMint(ctx); err != nil {
			log.Error().Err(err).Msg("监听safeMint失败")
		}
	}()

	// 初始化auction链监听器
	auctionListener, err := NFTAuction.NewListener(&cfg.Blockchain)
	if err != nil {
		log.Fatal().Err(err).Msg("初始化区块链监听器失败")
	}

	go func() {
		log.Info().Msg("启动NFTAuction监听器")
		if err := auctionListener.Start(ctx); err != nil {
			log.Error().Err(err).Msg("区块链监听器退出")
		}
	}()

	// 9. 监听退出信号（Ctrl+C、kill）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("开始关闭服务...")

	// 10. 优雅关闭
	cancel() // 停止监听器
	log.Info().Msg("服务已关闭")

}
