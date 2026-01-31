package redis

import (
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/ydh2333/NFTAuction-project/internal/models"
)

// Redis Key常量（统一管理，避免硬编码）
const (
	KeyAuctionHotRank = "nft_auction:hot_rank" // 拍卖热度排行Key
)

// IncrAuctionHot 增加拍卖热度（出价成功后调用，热度+1）
// auctionID: 拍卖ID（uint64）
func IncrAuctionHot(auctionID uint64) error {
	// ZINCRBY：有序集合中指定member的score+1
	return rdb.ZIncrBy(ctx, KeyAuctionHotRank, 1, strconv.FormatUint(auctionID, 10)).Err()
}

// GetTop5HotAuctions 获取Top5热门拍卖的拍卖ID（按热度降序，热度相同按最新出价时间降序）
// 返回：拍卖ID切片（从高到低）、错误
func GetTop5HotAuctions() ([]uint64, error) {
	// ZREVRANGE：按score降序取前5个member（0-4），不返回score
	res, err := rdb.ZRevRange(ctx, KeyAuctionHotRank, 0, 4).Result()
	if err != nil {
		return nil, err
	}

	// 字符串转uint64
	auctionIDs := make([]uint64, 0, len(res))
	for _, s := range res {
		aid, _ := strconv.ParseUint(s, 10, 64)
		auctionIDs = append(auctionIDs, aid)
	}
	return auctionIDs, nil
}

// InitHotRankFromDB 从MySQL初始化Redis热度排行（项目启动时调用）
// bids: 所有出价记录
func InitHotRankFromDB(bids []models.Bid) error {
	// 1. 先删除原有排行（避免脏数据）
	if err := rdb.Del(ctx, KeyAuctionHotRank).Err(); err != nil {
		return err
	}

	// 2. 统计每个拍卖的出价次数（热度值）
	hotMap := make(map[uint64]int64)
	for _, bid := range bids {
		hotMap[bid.AuctionID]++
	}

	// 3. 批量写入Redis（ZADD，批量操作提升性能）
	zs := make([]redis.Z, 0, len(hotMap))
	for aid, score := range hotMap {
		zs = append(zs, redis.Z{
			Score:  float64(score),
			Member: strconv.FormatUint(aid, 10),
		})
	}

	return rdb.ZAdd(ctx, KeyAuctionHotRank, zs...).Err()
}

// DelAuctionHot 从热度排行中删除拍卖（拍卖结束/流拍后可选调用，按需）
func DelAuctionHot(auctionID uint64) error {
	return rdb.ZRem(ctx, KeyAuctionHotRank, strconv.FormatUint(auctionID, 10)).Err()
}
