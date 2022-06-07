package cosmosdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

type paymentRepository struct {
	db *mongo.Client
}

func NewPaymentRepository() domain.PaymentRepository {
	mongoDBConnectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/Payment?authSource=admin", os.Getenv("DB_USER"), os.Getenv("DB_PW"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	fmt.Println("DB string : " + mongoDBConnectionString)
	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Client create error %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("unable to initialize connection %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("unable to connect %v", err)
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
		log.Fatalf("failed to list payment(s) %v", err)
	}
	err = rs.All(ctx, &payments)
	if err != nil {
		log.Fatalf("failed to list payment(s) %v", err)
	}
	return payments, nil
}

func (p *paymentRepository) PutPayment(ctx context.Context, model *domain.Payment) (*domain.Payment, error) {
	collection := p.db.Database("Payment").Collection("Payment")
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	fmt.Println("added Payment", result.InsertedID)
	return model, nil
}

// 履歴の登録
func (p *paymentRepository) PutPaymentHistory(ctx context.Context, model *domain.Payment) error {
	collection := p.db.Database("Payment").Collection("PaymentAllocateHistory")
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return err
	}
	fmt.Println("added PaymentHistory", result.InsertedID)
	return nil
}

// AcceptNoに該当するクレジット決済情報を取得
func (p *paymentRepository) GetPayment(ctx context.Context, acceptNo string) (*domain.Payment, error) {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}

	collection := p.db.Database("Payment").Collection("Payment")
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed to list payment(s) %v", err)
	}
	var payments []domain.Payment
	err = rs.All(ctx, &payments)
	if err != nil {
		log.Fatalf("failed to list payment(s) %v", err)
	}
	if len(payments) == 0 {
		fmt.Println("no todos found")
		return nil, nil
	}
	if payments[0].DeleteFlag {
		return nil, nil
	}
	return &payments[0], nil
}

// acceptNoに該当する決済データ論理を削除
func (p *paymentRepository) DeletePayment(ctx context.Context, acceptNo string) error {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}
	update := bson.D{{Key: "DeleteFlag", Value: acceptNo}}
	collection := p.db.Database("Payment").Collection("Payment")
	result, err := collection.UpdateOne(ctx, filter, update)
	fmt.Println("added Payment", result.UpsertedID)
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
