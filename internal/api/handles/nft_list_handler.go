package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/ydh2333/NFTAuction-project/internal/service"
	"github.com/ydh2333/NFTAuction-project/utils"
)

type NFTListHandler struct {
	nftListService service.NFTListService
}

func NewNFTListHandler() *NFTListHandler {
	return &NFTListHandler{
		nftListService: service.NewNFTListService(),
	}
}

func (a *NFTListHandler) GetNFTList(c *gin.Context) {
	address := c.Param("address")
	nftList, err := a.nftListService.GetNFTList(address)
	if err != nil {
		utils.SendError(c, 500, "获取NFT列表失败")
		return
	}
	utils.SendSuccess(c, "获取NFT列表成功", nftList)
}
