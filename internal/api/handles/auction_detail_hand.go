package handles

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydh2333/NFTAuction-project/internal/service"
	"github.com/ydh2333/NFTAuction-project/utils"
)

type AuctionDetailHandler struct {
	auctiondetailService service.AuctionDetailService
}

func NewAuctionDetailHandler() *AuctionDetailHandler {
	return &AuctionDetailHandler{
		auctiondetailService: service.NewAuctionDetailService(),
	}
}

func (a *AuctionDetailHandler) GetAuctionDetail(c *gin.Context) {
	auctionId, _ := strconv.Atoi(c.Param("id"))
	bids, err := a.auctiondetailService.GetAuctionDetail(uint(auctionId))
	if err != nil {
		utils.SendError(c, 500, "获取拍卖详情失败")
		return
	}
	utils.SendSuccess(c, "获取拍卖详情成功", bids)
}
