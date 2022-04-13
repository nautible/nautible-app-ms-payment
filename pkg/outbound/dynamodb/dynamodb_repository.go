package dynamodb

import (
	"context"
	"fmt"
	"os"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const TABLE_PAYMENT string = "Payment"
const TABLE_PAYMENT_ALLOCATE_HISTORY string = "PaymentAllocateHistory"
const TABLE_SEQUENCE string = "Sequence"

type dynamoDb struct {
	db *dynamo.DB
}

func NewDynamoDbRepository() domain.DbRepository {
	db, err := createSession()
	if err != nil {
		panic(err)
	}
	return &dynamoDb{db: db}
}

// 決済データの登録
func (p *dynamoDb) PutPaymentItem(ctx context.Context, model *domain.PaymentModel) (*domain.PaymentModel, error) {
	table := p.db.Table(TABLE_PAYMENT)
	if err := table.Put(model).RunWithContext(ctx); err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	fmt.Println("accept : " + model.AcceptNo)
	return model, nil
}

// OrderNoに該当する決済データを取得
func (p *dynamoDb) GetPaymentItem(ctx context.Context, orderNo string) (*domain.PaymentModel, error) {
	table := p.db.Table(TABLE_PAYMENT)
	var result domain.PaymentModel
	if err := table.Get("OrderNo", orderNo).OneWithContext(ctx, &result); err != nil {
		return nil, err
	}

	if result.DeleteFlag {
		return nil, nil
	}
	return &result, nil
}

// orderNoに該当する決済データ論理を削除
func (p *dynamoDb) DeletePaymentItem(ctx context.Context, orderNo string) error {
	table := p.db.Table(TABLE_PAYMENT)

	var result domain.Payment
	return table.Update("OrderNo", orderNo).Set("DeleteFlag", true).ValueWithContext(ctx, &result)
}

// 履歴の登録
func (p *dynamoDb) PutPaymentHistory(ctx context.Context, model *domain.PaymentModel) error {
	table := p.db.Table(TABLE_PAYMENT_ALLOCATE_HISTORY)
	if err := table.Put(model).If("attribute_not_exists(RequestId)").RunWithContext(ctx); err != nil {
		return err
	}
	fmt.Println("accept RequestId : " + model.RequestId)
	return nil
}

// シーケンス取得
func (p *dynamoDb) Sequence(ctx context.Context) (*int, error) {
	var counter struct {
		Name           string
		SequenceNumber int
	}
	table := p.db.Table(TABLE_SEQUENCE)
	err := table.Update("Name", "Payment").Add("SequenceNumber", 1).ValueWithContext(ctx, &counter)
	if err != nil {
		return nil, err
	}
	return &counter.SequenceNumber, err

}

func createSession() (*dynamo.DB, error) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT")))
	return db, nil
}
