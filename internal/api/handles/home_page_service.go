package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/ydh2333/NFTAuction-project/internal/repository"
	"github.com/ydh2333/NFTAuction-project/internal/service"
	"github.com/ydh2333/NFTAuction-project/utils"
)

type HomePageHandler struct {
	homePageService service.HomePageService
}

func NewHomePageHandler() *HomePageHandler {
	return &HomePageHandler{
		homePageService: service.NewHomePageService(),
	}
}

func (h *HomePageHandler) PlatformStatistics(c *gin.Context) {

	auctionCount, bidCount := h.homePageService.PlatformStatistics()

	utils.SendSuccess(c, "统计数据获取成功", gin.H{
		"auctionCount": auctionCount,
		"bidCount":     bidCount,
	})
}

type AuctionRequestList struct {
	repository.AuctionSearchParams
	repository.SortParams
	utils.PageParams
}

func (h *HomePageHandler) SearchAuctionsList(c *gin.Context) {
	var req AuctionRequestList

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "参数绑定失败")
		return
	}

	AuctionDetails, err := h.homePageService.SearchAuctionsList(req.AuctionSearchParams, req.SortParams, req.PageParams)

	if err != nil {
		utils.SendError(c, 500, "搜索拍卖失败")
		return
	}

	utils.SendSuccess(c, "搜索拍卖成功", AuctionDetails)
}
