package cosmosdb

import (
	"context"
	"os"
	"strings"
	"time"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type creditRepository struct {
	db *mongo.Client
}

func NewCreditRepository() domain.CreditRepository {
	mongoDBConnectionString := os.Getenv("DB_CONNECTION_STRING")
	mongoDBConnectionString = strings.Replace(mongoDBConnectionString, "${DB_USER}", os.Getenv("DB_USER"), -1)
	mongoDBConnectionString = strings.Replace(mongoDBConnectionString, "${DB_PW}", os.Getenv("DB_PW"), -1)
	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		zap.S().Fatalw("Client create error : " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		zap.S().Fatalw("unable to initialize connection : " + err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		zap.S().Fatalw("unable to connect : " + err.Error())
	}

	return &creditRepository{db: client}
}

func (p *creditRepository) Close() {
	p.db.Disconnect(context.TODO())
}

func (p *creditRepository) PutCreditPayment(ctx context.Context, model *domain.CreditPayment) (*domain.CreditPayment, error) {
	collection := p.db.Database("Payment").Collection("CreditPayment")
	doc := bson.D{
		{Key: "AcceptNo", Value: model.AcceptNo},
		{Key: "AcceptDate", Value: model.AcceptDate},
		{Key: "OrderNo", Value: model.OrderNo},
		{Key: "OrderDate", Value: model.OrderDate},
		{Key: "CustomerId", Value: model.CustomerId},
		{Key: "TotalPrice", Value: model.TotalPrice},
		{Key: "DeleteFlag", Value: model.DeleteFlag},
	}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		zap.S().Fatalw("failed to put payment(s) : " + err.Error())
		return nil, err
	}
	if result.InsertedID != nil {
		zap.S().Infof("added CreditPayment ID : %v", result.InsertedID)
	}
	return model, nil
}

// AcceptNoに該当するクレジット決済情報を取得
func (p *creditRepository) GetCreditPayment(ctx context.Context, acceptNo string) (*domain.CreditPayment, error) {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}

	collection := p.db.Database("Payment").Collection("CreditPayment")
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		zap.S().Fatalw("failed to list payment(s) : " + err.Error())
	}
	var payments []domain.CreditPayment
	err = rs.All(ctx, &payments)
	if err != nil {
		zap.S().Fatalw("failed to list payment(s) : " + err.Error())
	}
	if len(payments) == 0 {
		zap.S().Infow("no todos found")
		return nil, nil
	}
	if payments[0].DeleteFlag {
		return nil, nil
	}
	return &payments[0], nil
}

// acceptNoに該当する決済データ論理を削除
func (p *creditRepository) DeleteCreditPayment(ctx context.Context, acceptNo string) error {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}
	update := bson.D{{Key: "DeleteFlag", Value: true}}
	collection := p.db.Database("Payment").Collection("CreditPayment")
	result, err := collection.UpdateOne(ctx, filter, update)
	if result.UpsertedID != nil {
		zap.S().Infof("deleted CreditPayment ID : %v", result.UpsertedID)
	}
	return err
}

// シーケンス取得
func (p *creditRepository) Sequence(ctx context.Context) (*int, error) {
	filter := bson.D{{Key: "_id", Value: "CreditPayment"}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "SequenceNumber", Value: 1}}}}
	ops := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	collection := p.db.Database("Common").Collection("Sequence")
	var updatedDoc bson.M
	collection.FindOneAndUpdate(ctx, filter, update, ops).Decode(&updatedDoc)
	seq := int(updatedDoc["SequenceNumber"].(int32))
	return &seq, nil
}
