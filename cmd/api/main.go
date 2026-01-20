package main

import (
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/utils/logger"
)

func main() {
	// 初始化日志
	logger.InitLogger()
	cfg := config.LoadConfig()

	repository.InitDB(&cfg.MySQL)
}
