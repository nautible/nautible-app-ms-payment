# 動作確認用テストデータ

## データ作成

```
curl http://localhost:8080/payment/create -D - -X POST -i -H "Content-Type: application/json" -d "{\"orderNo\": \"1111-2222-3333-4444\", \"orderDate\": \"2021/03/01\", \"customerId\": 1, \"totalPrice\": 1300, \"paymentType\", \"03\"}"
```

## データ更新

```
curl http://localhost:8080/payment/ -D - -X PUT -i -H "Content-Type: application/json" -d "{\"paymentNo\": \"1\", \"orderNo\": \"1111-2222-3333-4444\", \"orderDate\": \"2021/03/01\", \"customerId\": 1, \"totalPrice\": 2300, \"acceptNo\":\"<createで作成した値>\"}"
```

## 検索

なし

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
