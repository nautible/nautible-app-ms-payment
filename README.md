# nautible-app-ms-payment

nautible-app-ms-payment project

## アーキテクチャ図

![アーキテクチャイメージ](./assets/architecture.svg)

## 機能

- cash
  - 代引き決済を行うダミーサービス
    - 決済登録
    - 決済キャンセル（論理削除）
- credit
  - クレジット決済を行うダミーサービス
    - 決済登録
    - 決済キャンセル（論理削除）
- bff
  - 各種決済サービスを呼び出して結果を返す

## ディレクトリ構成

[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を参考に構成


## サンプルアプリ利用手順

### skaffoldによるアプリケーション起動

BFF

```bash
cd scripts/bff
./skaffold.sh
```

Cash

```bash
cd scripts/cash
./skaffold.sh
```

Credit

```bash
cd scripts/credit
./skaffold.sh
```

## アプリ構築手順

### Golangバージョン

1.18

### OpenAPI

- oapi-codegenを導入

```bash
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.0
```

- YAMLファイルを準備
  - 参考：api/内のYAMLファイル

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

```bash
go mod init github.com/nautible/nautible-app-ms-payment
go mod tidy
```
