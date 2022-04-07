package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// エラー応答用JSON
type ErrorMessage struct {
	Message string `json:"message"`
	Detail  []ErrorDetail
}

// エラー詳細
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaymentService interface {
	CreatePayment(context.Context, *Payment)
	GetByOrderNo(context.Context, string, string) (*Payment, error)
	DeleteByOrderNo(context.Context, string, string) error
}

type PaymentStruct struct {
	history *DynamoDbRepository
	cash    *CashRepository
	credit  *CreditRepository
	order   *OrderRepository
}

func NewPaymentService(history *DynamoDbRepository, cash *CashRepository, credit *CreditRepository, order *OrderRepository) PaymentService {
	return &PaymentStruct{history, cash, credit, order}
}

// バックエンドサービスに支払作成処理を投げ、結果をOrderに返す
func (svc *PaymentStruct) CreatePayment(ctx context.Context, payment *Payment) {
	fmt.Println("CreatePaymentService")
	// バリデート
	var orderResponse OrderResponse
	result := validate(payment)
	if result != "" {
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = payment.RequestId
		orderResponse.Message = result
		orderResponse.Status = http.StatusBadRequest
		(*svc.order).PaymentResponse(ctx, &orderResponse)
		fmt.Println(orderResponse.Message)
		return
	}

	// 冪等性担保 履歴テーブルへの登録
	if err := (*svc.history).PutPaymentHistory(ctx, payment); err != nil {
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = payment.RequestId
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			// エラー内容が登録済みの場合は正常応答
			orderResponse.Status = http.StatusOK
			orderResponse.Message = messsageFormat("Success")
		} else {
			// それ以外のエラーは異常応答
			orderResponse.Status = http.StatusInternalServerError
			orderResponse.Message = messsageFormat(err.Error())
		}
		(*svc.order).PaymentResponse(ctx, &orderResponse)
		return
	}

	// バックエンドの決済処理を呼び出す
	orderResponse = *createPayment(ctx, svc, payment)

	// レスポンスをOrderに送信
	(*svc.order).PaymentResponse(ctx, &orderResponse)
}

func (svc *PaymentStruct) GetByOrderNo(ctx context.Context, paymentType string, orderNo string) (*Payment, error) {
	if paymentType == string(TypeCash) {
		return (*svc.cash).GetByOrderNo(ctx, orderNo)
	} else if paymentType == string(TypeCredit) {
		return (*svc.credit).GetByOrderNo(ctx, orderNo)
	}
	return nil, nil
}

func (svc *PaymentStruct) DeleteByOrderNo(ctx context.Context, paymentType string, orderNo string) error {
	if paymentType == string(TypeCash) {
		return (*svc.cash).DeleteByOrderNo(ctx, orderNo)
	} else if paymentType == string(TypeCredit) {
		return (*svc.credit).DeleteByOrderNo(ctx, orderNo)
	}
	return nil
}

func validate(paymentItem *Payment) string {
	// validator doc https://github.com/go-playground/validator/tree/master
	validate := validator.New()
	// jsonタグ名で応答できるように設定する
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	errorMessage := ErrorMessage{
		Message: "Validation Error",
		Detail:  []ErrorDetail{},
	}
	if err := validate.Struct(paymentItem); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, errrs := range errors {
			var detail ErrorDetail
			detail.Field = errrs.Field()
			detail.Message = errrs.Error()
			errorMessage.Detail = append(errorMessage.Detail, detail)
		}
		e, err := json.Marshal(errorMessage)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		result := string(e)
		return result
	}
	return ""
}

func createPayment(ctx context.Context, svc *PaymentStruct, payment *Payment) *OrderResponse {
	var orderResponse OrderResponse
	var res *Payment
	var err error
	if payment.PaymentType == string(TypeCash) {
		res, err = (*svc.cash).CreatePayment(ctx, payment)
	} else if payment.PaymentType == string(TypeCredit) {
		res, err = (*svc.credit).CreatePayment(ctx, payment)
	} else {
		// 支払い区分不正
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = res.RequestId
		orderResponse.Status = http.StatusBadRequest
		orderResponse.Message = messsageFormat("支払い区分が不正です paymentType : " + payment.PaymentType)
		return &orderResponse
	}
	// エラー発生
	if err != nil {
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = payment.RequestId
		orderResponse.Status = http.StatusInternalServerError
		orderResponse.Message = messsageFormat(err.Error())
		return &orderResponse
	}
	// 正常応答
	orderResponse.ProcessType = string(TypePayment)
	orderResponse.RequestId = res.RequestId
	orderResponse.Status = http.StatusOK
	orderResponse.Message = messsageFormat("Success")
	return &orderResponse
}

func messsageFormat(message string) string {
	errorMessage := ErrorMessage{
		Message: message,
		Detail:  []ErrorDetail{},
	}
	e, err := json.Marshal(errorMessage)
	result := string(e)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return result
}
