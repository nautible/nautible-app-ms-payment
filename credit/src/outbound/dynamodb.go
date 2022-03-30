package outbound

import (
	"context"
	"fmt"
	"time"

	"payment-credit/src/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type paymentStruct struct {
	db *dynamo.DB
}

func NewPaymentDB() domain.PaymentRepository {
	db, err := createSession()
	if err != nil {
		panic(err)
	}
	return &paymentStruct{db: db}
}

// 決済データの登録
func (p *paymentStruct) PutItem(ctx context.Context, model *domain.PaymentItem) (*domain.PaymentItem, error) {
	paymentNo, err := sequence(ctx, p.db)
	if err != nil {
		return nil, err
	}
	model.PaymentNo = fmt.Sprintf("C%09d", *paymentNo) // dummy クレジットの支払い番号はC始まりとする
	model.AcceptNo = fmt.Sprintf("A%09d", *paymentNo)  // dummy 受付番号はA始まりとする
	model.ReceiptDate = time.Now().String()            // dummy
	table := p.db.Table("Payment")
	if err := table.Put(model).Run(); err != nil {
		fmt.Printf("Failed to put item[%v]\n", err)
		return nil, err
	}
	fmt.Println("accept : " + model.AcceptNo)
	return model, nil
}

// paymentNoに該当する決済データを取得
func (p *paymentStruct) GetItem(ctx context.Context, paymentNo string) (*domain.PaymentItem, error) {
	table := p.db.Table("Payment")
	var result domain.PaymentItem
	err := table.Get("PaymentNo", paymentNo).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func createSession() (*dynamo.DB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://payment-credit-localstack.nautible-app-ms.svc.cluster.local:4566"),
		Credentials: credentials.NewStaticCredentials("test-key", "test-secret", ""),
	})
	if err != nil {
		return nil, err
	}
	db := dynamo.New(sess)
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
