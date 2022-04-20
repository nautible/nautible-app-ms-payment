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

type creditRepository struct {
	db *dynamo.DB
}

func NewCreditRepository() domain.CreditRepository {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT")))
	return &creditRepository{db: db}
}

// 決済データの登録
func (p *creditRepository) PutCreditPayment(ctx context.Context, model *domain.CreditPayment) (*domain.CreditPayment, error) {
	table := p.db.Table("CreditPayment")
	if err := table.Put(model).RunWithContext(ctx); err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	return model, nil
}

// AcceptNoに該当するクレジット決済情報を取得
func (p *creditRepository) GetCreditPayment(ctx context.Context, acceptNo string) (*domain.CreditPayment, error) {
	table := p.db.Table("CreditPayment")
	var result domain.CreditPayment
	if err := table.Get("AcceptNo", acceptNo).OneWithContext(ctx, &result); err != nil {
		return nil, err
	}
	if result.DeleteFlag {
		return nil, nil
	}
	return &result, nil
}

// acceptNoに該当する決済データ論理を削除
func (p *creditRepository) DeleteCreditPayment(ctx context.Context, acceptNo string) error {
	table := p.db.Table("CreditPayment")

	var result domain.Payment
	return table.Update("AcceptNo", acceptNo).Set("DeleteFlag", true).ValueWithContext(ctx, &result)
}

// シーケンス取得
func (p *creditRepository) Sequence(ctx context.Context) (*int, error) {
	var counter struct {
		Name           string
		SequenceNumber int
	}
	table := p.db.Table("Sequence")
	err := table.Update("Name", "CreditPayment").Add("SequenceNumber", 1).ValueWithContext(ctx, &counter)
	if err != nil {
		return nil, err
	}
	return &counter.SequenceNumber, err

}
