# 動作確認用テストデータ

# Cash APP

## データ作成

```
$ curl http://localhost:8080/payment/ -D - -X POST -i -H "Content-Type: application/json" -d "{\"orderNo\": \"1111-2222-3333-4444\", \"orderDate\": \"2021/03/01\", \"customerId\": 1, \"totalPrice\": 1300}"
```

## データ更新

```
curl http://localhost:8080/payment/ -D - -X PUT -i -H "Content-Type: application/json" -d "{\"orderNo\": 1, \"orderDate\": \"2021/03/01\", \"customerId\": 1, \"totalPrice\": 1300, \"productPrice\": 1000, \"tax\": 100, \"deliveryFee\": 200, \"orderStatus\": \"PROCESS\", \"orderDetail\": [{\"id\": 1, \"price\": 400, \"count\": 1},{\"id\": 2, \"price\": 600, \"count\": 1}], \"payment\": {\"paymentId\": \"DAIBIKI\"}, \"destination\": {\"zipCode\": \"583-0017\", \"address1\": \"osaka\", \"address2\": \"senri\", \"tel\": \"111-1111-1111\"}}"
```

## 検索

```
curl "http://localhost:8080/payment?customerId=1&orderDateFrom=2021-01-01&orderDateTo=2023-04-01" -D -
```

## ID指定

```
curl http://localhost:8080/payment/1 -D -
```

## データ削除

```
curl http://localhost:8080/payment/1 -X DELETE -D -
```


## local-stackのDynamoDBデータ確認

```
aws dynamodb scan --table-name Payment --endpoint-url http://localhost:4566
```
