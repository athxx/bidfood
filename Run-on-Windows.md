# How to run

This document provides instructions on how to run the BidFood application in **Windows** OS.

## Install dependencies

open cmd and run blow commands

```sh
## protoc support
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
## grpc support
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
## grpc test support
go install github.com/grpc-ecosystem/grpcurl/cmd/grpcurl@latest
```

## Run api

```sh
cd bidapi
go run ./...
```

## Run rpc

```sh
cd bidrpc
go run ./...
```

## Run unit tests

```sh
go test ./...
```

## Run api tests

**before you run this test, you must run the `api` and `rpc` servers**

1. install vscode plugin `REST Client` and open `test.http` this file
   ![2](/images/1.png)

1. click `Send Request` button
   ![2](/images/2.png)
   ![3](/images/3.png)

## Run grpc test

open `cmd` and run blow commands

```sh
# List Products
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{}' localhost:9000 bidrpcproto.ProductService/ListProducts

# ListProducts (Pagination)
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{"page":1,"page_size":2}' localhost:9000 bidrpcproto.ProductService/ListProducts

# CreateProduct
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{"name":"iPhone 15 Pro","description":"Latest iPhone","price":999.99,"quantity":50}' localhost:9000 bidrpcproto.ProductService/CreateProduct

# Get Single Product, please replace the product id
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{"id":"6584f023-9cfd-4fe0-b613-2b2ecf00fa4c"}' localhost:9000 bidrpcproto.ProductService/GetProduct

# UpdateProduct
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{"id":1,"name":"iPhone 15 Pro Max","description":"Updated desc","price":1099.99,"quantity":30}' localhost:9000 bidrpcproto.ProductService/UpdateProduct

# DeleteProduct
grpcurl -plaintext -proto "./bidrpc/bidrpcproto/product.proto" -d '{"id":"6584f023-9cfd-4fe0-b613-2b2ecf00fa4c"}' localhost:9000 bidrpcproto.ProductService/DeleteProduct

```
