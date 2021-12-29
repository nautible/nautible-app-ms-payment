
## golang


```
$ go mod init
$ go mod tidy
```

# cash

## OpenAPI

### install

```
$ go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```

支払いAPI

```
$ cd cash
$ midir generate
$ oapi-codegen -package generate -generate "types" -o src/generate/types.go openapi\payment.yaml
$ oapi-codegen -package generate -generate "spec" -o src/generate/spec.go openapi\payment.yaml
$ oapi-codegen -package generate -generate "chi-server" -o src/generate/server.go openapi\payment.yaml
```