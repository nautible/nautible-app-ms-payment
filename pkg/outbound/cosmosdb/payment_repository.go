package cosmosdb

import (
	"context"
	"os"
	"strings"
	"time"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type paymentRepository struct {
	db *mongo.Client
}

func NewPaymentRepository() domain.PaymentRepository {
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

	return &paymentRepository{db: client}
}

func (p *paymentRepository) Close() {
	p.db.Disconnect(context.TODO())
}

func (p *paymentRepository) FindPayment(ctx context.Context, customerId int32, orderDateFrom string, orderDateTo string) ([]*domain.Payment, error) {
	filter := bson.D{{Key: "CustomerId", Value: customerId}, {Key: "OrderDate", Value: bson.D{{Key: "$gt", Value: orderDateFrom}, {Key: "$lt", Value: orderDateTo}}}}
	var payments []*domain.Payment
	collection := p.db.Database("Payment").Collection("Payment")

	rs, err := collection.Find(ctx, filter)
	if err != nil {
		zap.S().Fatalw("failed to list payment(s) : " + err.Error())
	}
	err = rs.All(ctx, &payments)
	if err != nil {
		zap.S().Fatalw("failed to list payment(s) : " + err.Error())
	}
	return payments, nil
}

func (p *paymentRepository) PutPayment(ctx context.Context, model *domain.Payment) (*domain.Payment, error) {
	collection := p.db.Database("Payment").Collection("Payment")
	doc := bson.D{
		{Key: "RequestId", Value: model.RequestId},
		{Key: "PaymentNo", Value: model.PaymentNo},
		{Key: "PaymentType", Value: model.PaymentType},
		{Key: "OrderNo", Value: model.OrderNo},
		{Key: "OrderDate", Value: model.OrderDate},
		{Key: "AcceptNo", Value: model.AcceptNo},
		{Key: "AcceptDate", Value: model.AcceptDate},
		{Key: "CustomerId", Value: model.CustomerId},
		{Key: "TotalPrice", Value: model.TotalPrice},
		{Key: "OrderStatus", Value: model.OrderStatus},
		{Key: "DeleteFlag", Value: model.DeleteFlag},
	}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		zap.S().Errorw("Failed to put item : " + err.Error())
		return nil, err
	}
	if result.InsertedID != nil {
		zap.S().Infow("added Payment : " + result.InsertedID.(string))
	}
	return model, nil
}

// 履歴の登録
func (p *paymentRepository) PutPaymentHistory(ctx context.Context, model *domain.Payment) error {
	collection := p.db.Database("Payment").Collection("PaymentAllocateHistory")
	doc := bson.D{
		{Key: "RequestId", Value: model.RequestId},
		{Key: "PaymentNo", Value: model.PaymentNo},
		{Key: "PaymentType", Value: model.PaymentType},
		{Key: "OrderNo", Value: model.OrderNo},
		{Key: "OrderDate", Value: model.OrderDate},
		{Key: "AcceptNo", Value: model.AcceptNo},
		{Key: "AcceptDate", Value: model.AcceptDate},
		{Key: "CustomerId", Value: model.CustomerId},
		{Key: "TotalPrice", Value: model.TotalPrice},
		{Key: "OrderStatus", Value: model.OrderStatus},
		{Key: "DeleteFlag", Value: model.DeleteFlag},
	}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		zap.S().Errorw("Failed to put item : " + err.Error())
		return err
	}
	if result.InsertedID != nil {
		zap.S().Infow("added PaymentHistory : " + result.InsertedID.(string))
	}
	return nil
}

// orderNoに該当するクレジット決済情報を取得
func (p *paymentRepository) GetPayment(ctx context.Context, orderNo string) (*domain.Payment, error) {
	filter := bson.D{{Key: "OrderNo", Value: orderNo}}

	collection := p.db.Database("Payment").Collection("Payment")
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		zap.S().Fatalw("failed to list payment(s) : " + err.Error())
	}
	var payments []domain.Payment
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

// orderNoに該当する決済データ論理を削除
func (p *paymentRepository) DeletePayment(ctx context.Context, orderNo string) error {
	filter := bson.D{{Key: "OrderNo", Value: orderNo}}
	update := bson.D{{Key: "DeleteFlag", Value: true}}
	collection := p.db.Database("Payment").Collection("Payment")
	result, err := collection.UpdateOne(ctx, filter, update)
	if result.UpsertedID != nil {
		zap.S().Infow("added Payment : " + result.UpsertedID.(string))
	}
	return err
}

// シーケンス取得
func (p *paymentRepository) Sequence(ctx context.Context) (*int, error) {
	filter := bson.D{{Key: "_id", Value: "Payment"}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "SequenceNumber", Value: 1}}}}
	ops := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	collection := p.db.Database("Common").Collection("Sequence")
	var updatedDoc bson.M
	collection.FindOneAndUpdate(ctx, filter, update, ops).Decode(&updatedDoc)
	seq := int(updatedDoc["SequenceNumber"].(int32))
	return &seq, nil

}
