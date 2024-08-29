package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	cancelFunc context.CancelFunc
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	ctx, cancel := context.WithCancel(context.Background())
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
		cancelFunc: cancel,
	}
	go repo.MonitorAndCloseExpiredAuctions(ctx)
	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionDuration() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 10
	}

	return duration
}

func (ar *AuctionRepository) MonitorAndCloseExpiredAuctions(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Info("Checking for expired auctions")
			ar.checkAndCloseExpiredAuctions(ctx)
		case <-ctx.Done():
			logger.Info("Stopping MonitorAndCloseExpiredAuctions goroutine")
			return
		}
	}
}

func (ar *AuctionRepository) checkAndCloseExpiredAuctions(ctx context.Context) {
	auctionDuration := getAuctionDuration()
	now := time.Now().Unix()

	filter := ar.createExpiredAuctionsFilter(auctionDuration, now)
	update := ar.createUpdateCompletedStatus()

	ar.updateExpiredAuctions(ctx, filter, update)
}

type expiredAuctionsFilter struct {
	Timestamp primitive.M                  `bson:"timestamp"`
	Status    auction_entity.AuctionStatus `bson:"status"`
}

func (ar *AuctionRepository) createExpiredAuctionsFilter(auctionDuration time.Duration, currentTime int64) expiredAuctionsFilter {
	return expiredAuctionsFilter{
		Status:    auction_entity.Active,
		Timestamp: primitive.M{"$lt": currentTime - int64(auctionDuration.Seconds())},
	}
}

type updateCompletedStatus struct {
	SetStatus auction_entity.AuctionStatus `bson:"$set.status"`
}

func (ar *AuctionRepository) createUpdateCompletedStatus() updateCompletedStatus {
	return updateCompletedStatus{
		SetStatus: auction_entity.Completed,
	}
}

func (ar *AuctionRepository) updateExpiredAuctions(ctx context.Context, filter expiredAuctionsFilter, update updateCompletedStatus) {
	_, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update expired auctions", err)
	}
}

func (ar *AuctionRepository) Stop() {
	ar.cancelFunc()
}
