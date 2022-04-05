package domain

import (
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
	CreatePayment(*Payment)
	GetByOrderNo(string, string) (*Payment, error)
	DeleteByOrderNo(string, string) error
}

type PaymentStruct struct {
	cash   *CashRepository
	credit *CreditRepository
	order  *OrderRepository
}

func NewPaymentService(cash *CashRepository, credit *CreditRepository, order *OrderRepository) PaymentService {
	return &PaymentStruct{cash, credit, order}
}

// バックエンドサービスに支払作成処理を投げ、結果をOrderに返す
func (svc *PaymentStruct) CreatePayment(payement *Payment) {
	fmt.Println("CreatePaymentService")
	// バリデート
	var orderResponse OrderResponse
	result := validate(payement)
	if result != "" {
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = payement.RequestId
		orderResponse.Message = result
		orderResponse.Status = http.StatusBadRequest
		(*svc.order).PaymentResponse(&orderResponse)
		fmt.Println(orderResponse.Message)
		return
	}

	// バックエンドの決済処理を呼び出す
	var res *Payment
	var err error
	if payement.PaymentType == string(TypeCash) {
		res, err = (*svc.cash).CreatePayment(payement)
	} else if payement.PaymentType == string(TypeCredit) {
		res, err = (*svc.credit).CreatePayment(payement)
	} else {
		// 支払い区分不正
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = res.RequestId
		orderResponse.Status = http.StatusBadRequest
		orderResponse.Message = messsageFormat("支払い区分が不正です paymentType : " + payement.PaymentType)
		(*svc.order).PaymentResponse(&orderResponse)
		return
	}
	// エラー発生
	if err != nil {
		orderResponse.ProcessType = string(TypePayment)
		orderResponse.RequestId = payement.RequestId
		orderResponse.Status = http.StatusInternalServerError
		orderResponse.Message = messsageFormat(err.Error())
		(*svc.order).PaymentResponse(&orderResponse)
		return
	}
	// 正常応答
	orderResponse.ProcessType = string(TypePayment)
	orderResponse.RequestId = res.RequestId
	orderResponse.Status = http.StatusCreated
	orderResponse.Message = messsageFormat("Success")
	(*svc.order).PaymentResponse(&orderResponse)
}

func (svc *PaymentStruct) GetByOrderNo(paymentType string, orderNo string) (*Payment, error) {
	if paymentType == string(TypeCash) {
		return (*svc.cash).GetByOrderNo(orderNo)
	} else if paymentType == string(TypeCredit) {
		return (*svc.credit).GetByOrderNo(orderNo)
	}
	return nil, nil
}

func (svc *PaymentStruct) DeleteByOrderNo(paymentType string, orderNo string) error {
	if paymentType == string(TypeCash) {
		return (*svc.cash).DeleteByOrderNo(orderNo)
	} else if paymentType == string(TypeCredit) {
		return (*svc.credit).DeleteByOrderNo(orderNo)
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
