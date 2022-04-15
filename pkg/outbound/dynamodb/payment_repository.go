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

type dynamoDb struct {
	db *dynamo.DB
}

func NewPaymentRepository() domain.PaymentRepository {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT")))
	return &dynamoDb{db: db}
}

func (p *dynamoDb) FindPaymentItem(ctx context.Context, customerId int32, orderDateFrom string, orderDateTo string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	table := p.db.Table("Payment")

	if err := table.Get("CustomerId", customerId).Range("OrderDate", dynamo.GreaterOrEqual, orderDateFrom).Range("OrderDate", dynamo.LessOrEqual, orderDateTo).Index("GSI-CustomerId").AllWithContext(ctx, &payments); err != nil {
		fmt.Printf("Failed to get item[%v]\n", err)
		return nil, err
	}
	return payments, nil
}

// 決済データの登録
func (p *dynamoDb) PutPaymentItem(ctx context.Context, model *domain.Payment) (*domain.Payment, error) {
	table := p.db.Table("Payment")
	if err := table.Put(model).RunWithContext(ctx); err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	return model, nil
}

// OrderNoに該当する決済データを取得
func (p *dynamoDb) GetPaymentItem(ctx context.Context, orderNo string) (*domain.Payment, error) {
	table := p.db.Table("Payment")
	var result domain.Payment
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
	table := p.db.Table("Payment")

	var result domain.Payment
	return table.Update("OrderNo", orderNo).Set("DeleteFlag", true).ValueWithContext(ctx, &result)
}

// 履歴の登録
func (p *dynamoDb) PutPaymentHistory(ctx context.Context, model *domain.Payment) error {
	table := p.db.Table("PaymentAllocateHistory")
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
	table := p.db.Table("Sequence")
	err := table.Update("Name", "Payment").Add("SequenceNumber", 1).ValueWithContext(ctx, &counter)
	if err != nil {
		return nil, err
	}
	return &counter.SequenceNumber, err

}
