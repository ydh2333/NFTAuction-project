package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/ydh2333/NFTAuction-project/config"
)

// 全局Redis客户端
var rdb *redis.Client
var ctx = context.Background()

// 初始化Redis连接
func Init(cfg *config.RedisConfig) {
	rdb = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 5, // 最小空闲连接数
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Redis连接失败")
	}

	log.Info().Str("addr", cfg.Addr).Int("db", cfg.DB).Msg("Redis初始化成功")
}

// 关闭Redis连接（优雅关闭用）
func CloseRedis() {
	if err := rdb.Close(); err != nil {
		log.Error().Err(err).Msg("Redis关闭失败")
	}
	log.Info().Msg("Redis已关闭")
}
