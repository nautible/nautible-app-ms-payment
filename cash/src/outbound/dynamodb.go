package outbound

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"cash/src/domain"
)

type paymentStruct struct {
	db *dynamodb.DynamoDB
}

func NewPaymentDB(db *dynamodb.DynamoDB) domain.PaymentRepository {
	return &paymentStruct{db: db}
}
func (p *paymentStruct) PutItem(ctx context.Context, model domain.PaymentItem) error {

	av, err := dynamodbattribute.MarshalMap(model)
	if err != nil {
		panic(err)
	}
	_, result := p.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("Payment"),
		Item:      av,
	})
	return result
}

func (p *paymentStruct) GetItem(ctx context.Context, paymentNo string) domain.PaymentItem {

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://payment-localstack:4566"),
		Credentials: credentials.NewStaticCredentials("test-key", "test-secret", ""),
	})
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(session)

	// keyCond := expression.Key("DeviceID").Equal(expression.Value(deviceID)).
	// 	And(expression.Key("Timestamp").Between(
	// 			expression.Value(start.Format(time.RFC3339)),
	// 			expression.Value(end.Format(time.RFC3339))))
	keyCond := expression.Key("PaymentNo").Equal(expression.Value(paymentNo))

	// filterCond := expression.Name("DeviceType").Equal(expression.Value("Normal")).
	// 	And(expression.Name("CreatedYear").GreaterThan(expression.Value(2018)))

	// expr, err := expression.NewBuilder().WithKeyCondition(keyCond).WithFilter(filterCond).Build()
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		panic(err)
	}

	result, err := db.QueryWithContext(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("Payment"),
	})
	if err != nil {
		panic(err)
	}
	items := []domain.PaymentItem{}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	return items[0]
}
