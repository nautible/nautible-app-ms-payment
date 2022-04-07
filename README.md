# nautible-app-ms-payment

nautible-app-ms-payment project

## アーキテクチャ図

![アーキテクチャイメージ](./assets/architecture.svg)

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

## ディレクトリ構成

[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を参考に構成


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
oapi-codegen -package orderclient -generate "types" -o pkg/generate/orderclient/payment_types.go api/order.yaml
oapi-codegen -package orderclient -generate "client" -o pkg/generate/orderclient/http_client.go api/order.yaml
```

- 外部からの接続用サーバーコード生成

```bash
oapi-codegen -package bffserver -generate "types" -o pkg/generate/bffserver/types.go api/payment_bff.yaml
```

- BFFから内部API（cash/credit）へ接続するクライアントコード生成

```bash
oapi-codegen -package paymentclient -generate "types" -o pkg/generate/paymentclient/payment_types.go api/payment_backend.yaml
oapi-codegen -package paymentclient -generate "client" -o pkg/generate/paymentclient/http_client.go api/payment_backend.yaml
```

- 内部接続用（bffからcash/credit）のサーバーコード生成

```bash
oapi-codegen -package backendserver -generate "types" -o pkg/generate/backendserver/types.go api/payment_backend.yaml
oapi-codegen -package backendserver -generate "chi-server" -o pkg/generate/backendserver/server.go api/payment_backend.yaml
oapi-codegen -package backendserver -generate "spec" -o pkg/generate/backendserver/spec.go api/payment_backend.yaml
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
