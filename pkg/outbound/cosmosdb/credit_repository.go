package cosmosdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
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

	return &creditRepository{db: client}
}

func (p *creditRepository) Close() {
	p.db.Disconnect(context.TODO())
}

func (p *creditRepository) PutCreditPayment(ctx context.Context, model *domain.CreditPayment) (*domain.CreditPayment, error) {
	collection := p.db.Database("Payment").Collection("CreditPayment")
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	fmt.Println("added CreditPayment", result.InsertedID)
	return model, nil
}

// AcceptNoに該当するクレジット決済情報を取得
func (p *creditRepository) GetCreditPayment(ctx context.Context, acceptNo string) (*domain.CreditPayment, error) {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}

	collection := p.db.Database("Payment").Collection("CreditPayment")
	rs, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed to list payment(s) %v", err)
	}
	var payments []domain.CreditPayment
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
func (p *creditRepository) DeleteCreditPayment(ctx context.Context, acceptNo string) error {
	filter := bson.D{{Key: "AcceptNo", Value: acceptNo}}
	update := bson.D{{Key: "DeleteFlag", Value: acceptNo}}
	collection := p.db.Database("Payment").Collection("CreditPayment")
	result, err := collection.UpdateOne(ctx, filter, update)
	fmt.Println("added CreditPayment", result.UpsertedID)
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
