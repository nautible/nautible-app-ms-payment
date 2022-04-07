package outbound

import (
	"context"
	"fmt"
	"os"

	"github.com/nautible/nautible-app-ms-payment/pkg/cash/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type paymentStruct struct {
	db *dynamo.DB
}

func NewPaymentDB() domain.DynamoDbRepository {
	db, err := createSession()
	if err != nil {
		panic(err)
	}
	return &paymentStruct{db: db}
}

// 決済データの登録
func (p *paymentStruct) PutPaymentItem(ctx context.Context, model *domain.Payment) (*domain.Payment, error) {
	paymentNo, err := sequence(ctx, p.db)
	if err != nil {
		return nil, err
	}
	model.PaymentNo = fmt.Sprintf("P%10d", *paymentNo) // dummy 支払い番号はP始まりとする
	model.DeleteFlag = false
	table := p.db.Table("Payment")
	if err := table.Put(model).RunWithContext(ctx); err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	fmt.Println("accept : " + model.AcceptNo)
	return model, nil
}

// OrderNoに該当する決済データを取得
func (p *paymentStruct) GetPaymentItem(ctx context.Context, orderNo string) (*domain.Payment, error) {
	table := p.db.Table("Payment")
	var result domain.Payment
	err := table.Get("OrderNo", orderNo).OneWithContext(ctx, &result)
	if err != nil {
		return nil, err
	}
	if result.DeleteFlag {
		return nil, nil
	}
	return &result, nil
}

// orderNoに該当する決済データ論理を削除
func (p *paymentStruct) DeletePaymentItem(ctx context.Context, orderNo string) error {
	table := p.db.Table("Payment")

	var result domain.Payment
	return table.Update("OrderNo", orderNo).Set("DeleteFlag", true).ValueWithContext(ctx, &result)
}

func createSession() (*dynamo.DB, error) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, aws.NewConfig().WithRegion(os.Getenv("DYNAMODB_REGION")).WithEndpoint(os.Getenv("DYNAMODB_ENDPOINT")))
	return db, nil
}

// シーケンス取得
func sequence(ctx context.Context, db *dynamo.DB) (*int, error) {
	var counter struct {
		Name           string
		SequenceNumber int
	}
	table := db.Table("Sequence")
	err := table.Update("Name", "Payment").Add("SequenceNumber", 1).ValueWithContext(ctx, &counter)
	if err != nil {
		return nil, err
	}
	return &counter.SequenceNumber, err

}
