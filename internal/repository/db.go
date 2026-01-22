package repository

import (
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
	"github.com/ydh2333/NFTAuction-project/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.MySQLConfig) {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志（开发环境）
	})
	if err != nil {
		log.Fatal().Err(err).Msg("数据库连接失败")
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("获取数据库连接池失败")
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&models.NFT{},
		&models.Auction{},
		&models.Bid{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("数据库表迁移失败")
	}

	log.Info().Msg("数据库连接成功")
}
