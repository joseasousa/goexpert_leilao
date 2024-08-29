package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"
)

type AuctionRepositoryMock struct {
	Auctions map[string]*auction_entity.Auction
}

func NewAuctionRepositoryMock() *AuctionRepositoryMock {
	return &AuctionRepositoryMock{
		Auctions: make(map[string]*auction_entity.Auction),
	}
}

func (m *AuctionRepositoryMock) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	m.Auctions[auctionEntity.Id] = auctionEntity
	return nil
}

func (m *AuctionRepositoryMock) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	auction, ok := m.Auctions[id]
	if !ok {
		return nil, internal_error.NewInternalServerError("Auction not found")
	}
	return auction, nil
}

func (m *AuctionRepositoryMock) SaveAuction(auction *auction_entity.Auction) {
	m.Auctions[auction.Id] = auction
}

func (m *AuctionRepositoryMock) MonitorAndCloseExpiredAuctions(ctx context.Context) {
	now := time.Now().Unix()
	auctionDuration := getAuctionDuration()

	for id, auction := range m.Auctions {
		if auction.Status == auction_entity.Active && auction.Timestamp.Unix() < now-int64(auctionDuration.Seconds()) {
			auction.Status = auction_entity.Completed
			m.Auctions[id] = auction
		}
	}
}
