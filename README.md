# nautible-app-payment

nautible-app-payment project

# アーキテクチャ図

![アーキテクチャイメージ](./architecture.png)
# 機能

- cash
  - 代引き決済を行うダミーサービス
    - 代引き決済
    - 決済内容更新
    - 決済キャンセル
- convenience
  - コンビニ決済を行うダミーサービス
    - コンビニ決済（モック処理）
    - 決済内容更新
    - 決済キャンセル
- credit
  - クレジット決済を行うダミーサービス
    - クレジット決済（モック処理）
    - 決済内容更新
    - 決済キャンセル
- bff
  - 各種決済サービスを呼び出して結果を返すフロントサービス

# 以下TODO
# OpenAPI

REST通信のIF構築にはOpenAPIを使用している

## YAML

src/main/openapi/配下に定義ファイル(YAML)を用意

## Generate OpenAPI

```
$ mvn clean generate-sources
```

# ローカルで直接Quarkusアプリを起動

外部接続とかもなく、単体で動作確認するだけであればこれでよい

```
$ mvn compile quarkus:dev
```

# minikube

```
$ minikube start --docker-env http_proxy=<proxy> --docker-env https_proxy=<proxy> --docker-env no_proxy=localhost,158.201.0.0/16 --network-plugin=cni --enable-default-cni --driver=docker --kubernetes-version=v1.18.10
```

# skaffold

v1.20を使用（APIVersion：skaffold/v2beta12）

## namespace

minikube上にアプリケーション用のネームスペース作成

```
$ kubectl create namespace nautible-app
```

## skaffoldの起動

```
$ skaffold dev --port-forward 
```

# DynamoDB利用メモ

[Mapperの利用方法](https://docs.aws.amazon.com/ja_jp/sdk-for-java/latest/developer-guide/examples-dynamodb-enhanced.html)

[任意のオブジェクトからDynamoDBで保存できる形式へのコンバートはこちらを参考にコンバータを作る](https://github.com/aws/aws-sdk-java-v2/tree/master/services-custom/dynamodb-enhanced/src/main/java/software/amazon/awssdk/enhanced/dynamodb/internal/converter/attribute)

# Telepresence

## ローカルのDockerプロセスにイメージを登録

```
$ docker build -f src\main\docker\Dockerfile.fast-jar -t payment .
```

## ローカルにTelepresenceを準備(Windows)

```
$ git clone https://github.com/rashinban/rashinban-app-develop.git
$ cd rashinban-app-develop
$ build rashinban
```

## イメージの実行

```
$ docker run --net=host --ipc=host --uts=host --pid=host -it --security-opt=seccomp=unconfined --privileged --rm -v /tmp:/tmp -v /:/host -v /var/run/docker.sock:/var/run/docker.sock telepresence
```

## コンテナの起動

```
$ telepresence --method container --namespace <ネームスペース名> --new-deployment <デプロイメント名> --docker-run --rm <プロキシして動かすコンテナ>
```

デプロイメント名は任意の名前でよいが、自分のものとわかるようにしておくとよい

### すでにデプロイ済みのDeploymentを切り替える場合

```
$ telepresence --method container --namespace <ネームスペース名> --swap-deployment <切り替えるデプロイメント名> --docker-run --rm <プロキシして動かすコンテナのイメージ>
```

## [テストデータ](https://github.com/rashinban/nautible-app-payment/blob/main/testdata.md)

# 機能

- 注文情報登録
  - 出荷指示送信
- 注文情報変更
  - 出荷指示変更送信
- 注文情報削除
  - 出荷指示取り消し送信
- 注文情報取得（注文ID指定）
- 注文情報取得（顧客ID指定） 

削除内容いったんメモ

```
jobs:
  build:
    
    runs-on: ubuntu-latest

    env:
      IMAGE_TAG: ${{ github.sha }}

    steps:
    - name: Checkout repo
      uses: actions/checkout@v2
    - name: Checkout manifest repo
      uses: actions/checkout@v2
      with:
        repository: rashinban/nautible-app-payment-manifest
        path: nautible-app-payment-manifest
        #token: ${{ secrets.PAT }}
...
    - name: pull request
      id: pull-request
      env:
        #token: ${{ secrets.PAT }}
        tag: update-image-feature-${{ github.sha }}
      run: |
        cd $GITHUB_WORKSPACE/nautible-app-payment-manifest
        git checkout -b $tag
        sed -i 's/image: ${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.ap-northeast-1.amazonaws.com\/nautible-app-payment-bff:\(.*\)/image: ${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.ap-northeast-1.amazonaws.com\/nautible-app-payment-bff:'$IMAGE_TAG'/' ./bff/base/payment-bff-deploy.yaml 
        git config user.name github-actions
        git config user.email github-actions@github.com
        git add .
        git commit -m "update manifest"
        git push --set-upstream origin $tag
        curl -X POST -H "Accept: application/vnd.github.v3+json" -H "Authorization: token $token" "https://api.github.com/repos/rashinban/nautible-app-payment-manifest/pulls" -d '{"title": "new image deploy request", "head": "rashinban:'$tag'", "base": "main"}'
```