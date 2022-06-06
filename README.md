# nautible-app-ms-payment

## 概要

本プロジェクトは決済サンプルアプリケーションになります。アプリケーションはGolangで実装しており、以下の技術要素を含みます。

- oapi-codegenを利用したOpenAPIサーバー(chi)/クライアントの生成
- Dapr同期通信（ServiceInvocation）
  - oapi-codegenで生成したHTTPクライアントからのDapr同期通信（DaprSDKは未使用）
- Dapr非同期通信
  - net/httpパッケージで作成したHTTPサーバーでCloudEvents(application/octet-stream)の受信処理（DaprSDKは未使用）
- AWSSDKを利用したDynamoDBアクセス

## アーキテクチャ図

![アーキテクチャイメージ](./assets/architecture.svg)

## 機能

- credit
  - クレジット決済を行うダミーサービス
    - 決済登録
    - 決済キャンセル（論理削除）
- payment
  - 外部のサービスが決済サービスを呼び出すためのエンドポイント

なお、paymentとcreditはDaprのServiceInvocationの技術サンプルのため別プロセスで実行するようにしています。

## ディレクトリ構成

[Standard Go Project Layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を参考に構成

## サンプルアプリ利用手順

### skaffoldによるアプリケーション起動

Payment

```bash
cd scripts/payment
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

- PaymentからOrderサービスへ接続するクライアントコード生成

```bash
oapi-codegen -package orderclient -generate "types" -o pkg/generate/orderclient/payment_types.go api/order.yaml
oapi-codegen -package orderclient -generate "client" -o pkg/generate/orderclient/http_client.go api/order.yaml
```

- 外部からの接続用サーバーコード生成

```bash
oapi-codegen -package paymentserver -generate "types" -o pkg/generate/paymentserver/types.go api/payment.yaml
```

- Paymentから内部API（credit）へ接続するクライアントコード生成

```bash
oapi-codegen -package creditclient -generate "types" -o pkg/generate/creditclient/payment_types.go api/credit.yaml
oapi-codegen -package creditclient -generate "client" -o pkg/generate/creditclient/http_client.go api/credit.yaml
```

- 内部接続用（paymentからcredit）のサーバーコード生成

```bash
oapi-codegen -package creditserver -generate "types" -o pkg/generate/creditserver/types.go api/credit.yaml
oapi-codegen -package creditserver -generate "chi-server" -o pkg/generate/creditserver/server.go api/credit.yaml
oapi-codegen -package creditserver -generate "spec" -o pkg/generate/creditserver/spec.go api/credit.yaml
```

### go mod

```bash
go mod init github.com/nautible/nautible-app-ms-payment
go mod tidy
```

## ローカルでの実行

### payment

- aws用モジュールの実行

```bash
cd scripts/payment
./skaffold.sh aws
```

- azure用モジュールの実行

```bash
cd scripts/payment
./skaffold.sh azure
```

### credit

- aws用モジュールの実行

```bash
cd scripts/credit
./skaffold.sh aws
```

- azure用モジュールの実行

```bash
cd scripts/credit
./skaffold.sh azure
```
