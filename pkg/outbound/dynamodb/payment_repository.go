package dynamodb

import (
	"context"
	"os"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type PaymentRepository struct {
	db *dynamo.DB
}

func NewPaymentRepository() domain.PaymentRepository {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT")))
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) Close() {
	zap.S().Infow("DynamoDB close NoOp")
}

func (p *PaymentRepository) FindPayment(ctx context.Context, customerId int32, orderDateFrom string, orderDateTo string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	table := p.db.Table("Payment")

	if err := table.Get("CustomerId", customerId).Range("OrderDate", dynamo.GreaterOrEqual, orderDateFrom).Range("OrderDate", dynamo.LessOrEqual, orderDateTo).Index("GSI-CustomerId").AllWithContext(ctx, &payments); err != nil {
		zap.S().Errorw("Failed to get item : " + err.Error())
		return nil, err
	}
	return payments, nil
}

// 決済データの登録
func (p *PaymentRepository) PutPayment(ctx context.Context, model *domain.Payment) (*domain.Payment, error) {
	table := p.db.Table("Payment")
	if err := table.Put(model).RunWithContext(ctx); err != nil {
		zap.S().Errorw("Failed to put item : " + err.Error())
		return nil, err
	}
	return model, nil
}

// OrderNoに該当する決済データを取得
func (p *PaymentRepository) GetPayment(ctx context.Context, orderNo string) (*domain.Payment, error) {
	table := p.db.Table("Payment")
	var result domain.Payment
	if err := table.Get("OrderNo", orderNo).OneWithContext(ctx, &result); err != nil {
		zap.S().Errorw("Failed to get item : " + err.Error())
		return nil, err
	}
	if result.DeleteFlag {
		return nil, nil
	}
	return &result, nil
}

// orderNoに該当する決済データ論理を削除
func (p *PaymentRepository) DeletePayment(ctx context.Context, orderNo string) error {
	table := p.db.Table("Payment")

	var result domain.Payment
	return table.Update("OrderNo", orderNo).Set("DeleteFlag", true).ValueWithContext(ctx, &result)
}

// 履歴の登録
func (p *PaymentRepository) PutPaymentHistory(ctx context.Context, model *domain.Payment) error {
	table := p.db.Table("PaymentAllocateHistory")
	if err := table.Put(model).If("attribute_not_exists(RequestId)").RunWithContext(ctx); err != nil {
		zap.S().Errorw("Failed to put item : " + err.Error())
		return err
	}
	return nil
}

// シーケンス取得
func (p *PaymentRepository) Sequence(ctx context.Context) (*int, error) {
	var counter struct {
		Name           string
		SequenceNumber int
	}
	table := p.db.Table("Sequence")
	err := table.Update("Name", "Payment").Add("SequenceNumber", 1).ValueWithContext(ctx, &counter)
	if err != nil {
		zap.S().Errorw("Failed to update sequence : " + err.Error())
		return nil, err
	}
	return &counter.SequenceNumber, err

}
