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
	CreatePayment(*PaymentItem)
	GetByPaymentNo(string) (*PaymentItem, error)
}

type PaymentStruct struct {
	payment *PaymentRepository
	order   *OrderRepository
}

func NewPaymentService(payment *PaymentRepository, order *OrderRepository) PaymentService {
	return &PaymentStruct{payment, order}
}

// バックエンドサービスに支払作成処理を投げ、結果をOrderに返す
func (svc *PaymentStruct) CreatePayment(payementItem *PaymentItem) {
	fmt.Println("CreatePaymentService")
	// バリデート
	var orderResponse OrderResponse
	result := validate(payementItem)
	if result != "" {
		orderResponse.ProcessType = string(Payment)
		orderResponse.RequestId = payementItem.RequestId
		orderResponse.Message = result
		orderResponse.Status = http.StatusBadRequest
		(*svc.order).PaymentResponse(&orderResponse)
		fmt.Println(orderResponse.Message)
		return
	}

	//
	res, err := (*svc.payment).CreatePayment(payementItem)
	// エラー発生
	if err != nil {
		orderResponse.ProcessType = string(Payment)
		orderResponse.RequestId = res.RequestId
		orderResponse.Status = http.StatusInternalServerError
		orderResponse.Message = messsageFormat(err.Error())
		(*svc.order).PaymentResponse(&orderResponse)
		return
	}
	// 正常応答
	orderResponse.ProcessType = string(Payment)
	orderResponse.RequestId = res.RequestId
	orderResponse.Status = http.StatusCreated
	orderResponse.Message = messsageFormat("Success")
	(*svc.order).PaymentResponse(&orderResponse)
}

func (svc *PaymentStruct) GetByPaymentNo(paymentNo string) (*PaymentItem, error) {
	return (*svc.payment).GetByPaymentNo(paymentNo)
}

func validate(paymentItem *PaymentItem) string {
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
