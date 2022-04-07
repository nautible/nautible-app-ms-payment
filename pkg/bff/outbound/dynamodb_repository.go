package outbound

import (
	"context"
	"fmt"

	"github.com/nautible/nautible-app-ms-payment/pkg/bff/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://payment-bff-localstack.nautible-app-ms.svc.cluster.local:4566"),
		Credentials: credentials.NewStaticCredentials("test-key", "test-secret", ""),
	})
	if err != nil {
		return nil, err
	}
	db := dynamo.New(sess)
	return db, nil
}
