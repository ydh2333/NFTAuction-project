package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ydh2333/NFTAuction-project/internal/api/handles"
	middlewares "github.com/ydh2333/NFTAuction-project/internal/api/middleware"
	"github.com/ydh2333/NFTAuction-project/internal/blockchain"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine, blockchainClient *blockchain.AuctionContract) {
	r.Use(middlewares.Logger()) // 全局中间件
	r.Use(gin.Recovery())       // 异常恢复

	// 路由分组
	api := r.Group("/api")
	{
		// 首页相关接口
		homePageHandler := handles.NewHomePageHandler()
		homePage := api.Group("/hongPage")
		{
			homePage.GET("/platformStatistics", homePageHandler.PlatformStatistics)
			homePage.POST("/auctionList", homePageHandler.SearchAuctionsList)
		}

	}

}
