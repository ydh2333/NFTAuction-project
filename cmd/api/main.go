package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/api/routes"
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

	// 4. 创建全局上下文,统一管理所有监听器
	ctx, cancel := context.WithCancel(context.Background())
	// 确保程序退出时执行cancel，且只执行一次
	defer func() {
		cancel()
		log.Info().Msg("全局上下文已关闭，所有监听器停止")
	}()

	// 5. 初始化erc721SafeMint监听器
	erc721Listener, err := ERC721.NewERC721Listener(&cfg.Blockchain)
	if err != nil {
		log.Fatal().Err(err).Msg("初始化erc721SafeMint监听器失败")
	}

	// 启动ERC721监听器（后台协程，避免阻塞主线程）
	go func() {
		log.Info().Msg("启动ERC721 safeMint监听器")
		if err := erc721Listener.StartListeningSafeMint(ctx); err != nil {
			log.Error().Err(err).Msg("监听safeMint失败")
		}
	}()

	// 6. 初始化auction链监听器
	auctionListener, err := NFTAuction.NewListener(&cfg.Blockchain)
	if err != nil {
		log.Fatal().Err(err).Msg("初始化auction监听器失败")
	}

	go func() {
		log.Info().Msg("启动NFTAuction监听器")
		if err := auctionListener.Start(ctx); err != nil {
			log.Error().Err(err).Msg("NFTAuction监听器退出")
		}
	}()

	// 7. 初始化Gin
	gin.SetMode(gin.ReleaseMode) // 生产环境使用ReleaseMode
	r := gin.Default()

	// 8. 注册路由
	routes.InitRoutes(r)

	// 9. 启动HTTP服务
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 异步启动 HTTP 服务
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("HTTP服务启动成功，监听中...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP服务启动失败")
		}
	}()

	// 10. 监听退出信号（Ctrl+C、kill）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("开始关闭服务...")

	// 11. 优雅关闭
	srv.Shutdown(ctx)    // 关闭HTTP服务
	cancel()             // 停止所有区块链监听器（通过上下文取消）
	repository.CloseDB() // 关闭数据库
	defer log.Info().Msg("所有服务已关闭")
}
