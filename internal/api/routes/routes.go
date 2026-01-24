package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ydh2333/NFTAuction-project/internal/api/handles"
	middlewares "github.com/ydh2333/NFTAuction-project/internal/api/middleware"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
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
		// 拍卖详情页面
		auctionDetailHandler := handles.NewAuctionDetailHandler()
		auctionDetail := api.Group("/auctionDetail")
		{
			auctionDetail.GET("/:id", auctionDetailHandler.GetAuctionDetail)
		}
		// 个人主页/NFT列表
		nftListHandler := handles.NewNFTListHandler()
		nftList := api.Group("/ownerPage")
		{
			nftList.GET("/nftList/:address", nftListHandler.GetNFTList)
		}

	}

}
