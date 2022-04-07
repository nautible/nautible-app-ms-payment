package outbound

import (
	"context"
	"fmt"
	"os"

	"github.com/nautible/nautible-app-ms-payment/pkg/bff/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type paymentStruct struct {
	db *dynamo.DB
}

func NewDynamoDbRepository() domain.DynamoDbRepository {
	db, err := createSession()
	if err != nil {
		panic(err)
	}
	return &paymentStruct{db: db}
}

// 履歴の登録
func (p *paymentStruct) PutPaymentHistory(ctx context.Context, model *domain.Payment) error {
	table := p.db.Table("PaymentHistory")
	if err := table.Put(model).If("attribute_not_exists(RequestId)").RunWithContext(ctx); err != nil {
		return err
	}
	fmt.Println("accept RequestId : " + model.RequestId)
	return nil
}

func createSession() (*dynamo.DB, error) {
	sess := session.Must(session.NewSession())
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(endpoint))
	return db, nil
}
