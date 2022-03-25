# nautible-app-ms-payment

nautible-app-ms-payment project

## アーキテクチャ図

![アーキテクチャイメージ](./architecture.svg)

## 機能

- cash
  - 代引き決済を行うダミーサービス
    - 代引き決済
    - 決済内容更新
    - 決済キャンセル
- credit
  - クレジット決済を行うダミーサービス
    - クレジット決済（モック処理）
    - 決済内容更新
    - 決済キャンセル
- bff
  - 各種決済サービスを呼び出して結果を返すフロントサービス

## サンプルアプリ利用手順


## アプリ構築手順

### 前提

golang(v16)はインストール済みとする

### OpenAPI

- oapi-codegenを導入

```bash
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0
```

- YAMLファイルを準備
  - 参考：bff/openapi/内のYAMLファイル

- BFFからOrderサービスへ接続するクライアントコード生成

```bash
oapi-codegen -package order -generate "types" -o bff/generate/client/order/payment_types.go openapi/order.yaml
oapi-codegen -package order -generate "client" -o bff/generate/client/order/http_client.go openapi/order.yaml
```

- 外部からの接続用サーバーコード生成

```bash
oapi-codegen -package server -generate "types" -o bff/generate/server/types.go openapi/payment_bff.yaml
```

- BFFから内部API（cash/credit）へ接続するクライアントコード生成

```bash
oapi-codegen -package outbound -generate "types" -o bff/generate/client/payment/payment_types.go openapi/payment_backend.yaml
oapi-codegen -package outbound -generate "client" -o bff/generate/client/payment/http_client.go openapi/payment_backend.yaml
```

- 内部接続用（bffからcash/credit）のサーバーコード生成

```bash
oapi-codegen -package server -generate "types" -o cash/generate/server/types.go openapi/payment_backend.yaml
oapi-codegen -package server -generate "chi-server" -o cash/generate/server/server.go openapi/payment_backend.yaml
oapi-codegen -package server -generate "spec" -o cash/generate/server/spec.go openapi/payment_backend.yaml
oapi-codegen -package server -generate "types" -o credit/generate/server/types.go openapi/payment_backend.yaml
oapi-codegen -package server -generate "chi-server" -o credit/generate/server/server.go openapi/payment_backend.yaml
oapi-codegen -package server -generate "spec" -o credit/generate/server/spec.go openapi/payment_backend.yaml
```

### go mod

BFF

```bash
cd bff
go mod init payment-bff
go mod tidy
```

Cash

```bash
cd cash
go mod init payment-cash
go mod tidy
```

Credit

```bash
cd credit
go mod init payment-credit
go mod tidy
```
