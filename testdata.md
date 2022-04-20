# 動作確認用テストデータ

## 決済サービス

- 決済サービス（payment）を起動し、ポートフォワード
- クレジットサービス（credit）を起動（paymentType:01の場合、呼び出されるため）

```bash
cd scripts/payment
./skaffold.sh --port-forward
cd ../credit
./skaffold.sh
```

## データ作成

クレジット決済（決済サービスおよびクレジットサービスを呼び出すテスト）

```bash
curl -X POST http://localhost:8080/payment/create \
    -H "Content-Type: application/cloudevents+json" \
    -d '{"datacontenttype":"application/json","id":"1","source":"curl","type":"http request","specversion":"1.0", "data":"{\"requestId\":\"O0000000001\",\"orderDate\":\"2022-01-02T03:04:05\",\"customerId\":1,\"totalPrice\":10000,\"paymentType\":\"01\",\"orderNo\":\"O0000000001\"}"}'
```

現金代引き（決済サービスのみを呼び出すテスト）

```bash
curl -X POST http://localhost:8080/payment/create \
    -H "Content-Type: application/cloudevents+json" \
    -d '{"datacontenttype":"application/json","id":"1","source":"curl","type":"http request","specversion":"1.0", "data":"{\"requestId\":\"O0000000002\",\"orderDate\":\"2022-01-02T03:04:05\",\"customerId\":1,\"totalPrice\":10000,\"paymentType\":\"02\",\"orderNo\":\"O0000000002\"}"}'
```

## データ削除

```bash
curl -X POST http://localhost:8080/payment/rejectCreate \
    -H "Content-Type: application/cloudevents+json" \
    -d '{"datacontenttype":"application/json","id":"1","source":"curl","type":"http request","specversion":"1.0", "data":"{\"orderNo\":\"O0000000001\"}"}'
```

## クレジットデータ

- 決済サービスの起動（local-stackを起動するため）
- クレジットサービスを起動し、ポートフォワード

```bash
cd scripts/payment
./skaffold.sh
cd ../credit
./skaffold.sh --port-forward
```

### データの作成

データは"A"+10桁の連番で登録されている

```bash
 curl -X POST "http://localhost:8080/credit/A0000000001" -D -
```

### データの確認

データは"A"+10桁の連番で登録されている

```bash
 curl "http://localhost:8080/credit/A0000000001" -D -
```

### データの削除

```bash
 curl -X DELETE "http://localhost:8080/credit/A0000000001" -D -
```

## local-stackのDynamoDBデータ確認

```bash
kubectl port-forward svc/payment-localstack -n nautible-app-ms 4566:4566
aws dynamodb scan --table-name Payment --endpoint-url http://localhost:4566
aws dynamodb scan --table-name CreditPayment --endpoint-url http://localhost:4566
aws dynamodb scan --table-name Sequence --endpoint-url http://localhost:4566
```
