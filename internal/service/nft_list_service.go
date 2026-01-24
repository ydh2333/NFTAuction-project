package service

import (
	"github.com/ydh2333/NFTAuction-project/internal/repository"
)

type NFTListService interface {
	GetNFTList(OwnerAddress string) ([]repository.NftDetail, error)
}

type nftListService struct {
	nftRepo repository.NFTRepository
}

func NewNFTListService() NFTListService {
	return &nftListService{nftRepo: repository.NewNFTRepository()}
}

func (s *nftListService) GetNFTList(OwnerAddress string) ([]repository.NftDetail, error) {
	return s.nftRepo.GetNFTByOwnerAddress(OwnerAddress)
}
