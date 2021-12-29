
## OpenAPI

cash,convenience,creditサービスへ接続するためのクライアント生成

```
$ oapi-codegen -package payment -generate "client" -o src/generate/client/payment/client.go openapi\payment.yaml
$ oapi-codegen -package payment -generate "types" -o src/generate/client/payment/types.go openapi\payment.yaml
```

orderサービスから呼び出されるサーバー作成

```
$ oapi-codegen -package payment -generate "types" -o src/generate/server/payment/types.go openapi\payment.yaml
$ oapi-codegen -package payment -generate "spec" -o src/generate/server/payment/spec.go openapi\payment.yaml
$ oapi-codegen -package payment -generate "chi-server" -o src/generate/server/payment/server.go openapi\payment.yaml
```

orderサービスへ接続するためのクライアント生成

```
$ oapi-codegen -package order -generate "types" -o src/generate/server/order/client.go openapi\order.yaml
$ oapi-codegen -package order -generate "spec" -o src/generate/server/order/types.go openapi\order.yaml
```