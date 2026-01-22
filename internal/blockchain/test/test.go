package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/blockchain"
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

	// 4. 初始化区块链监听器
	listener, err := blockchain.NewListener(&cfg.Blockchain)
	if err != nil {
		// logger.Log.Fatal().Err(err).Msg("初始化区块链监听器失败")
		log.Error().Err(err).Msg("初始化区块链监听器失败")
	}

	// 5. 启动监听器（后台协程）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info().Msg("listener start")
		if err := listener.Start(ctx); err != nil {
			// logger.Log.Fatal().Err(err).Msg("区块链监听器退出")
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
