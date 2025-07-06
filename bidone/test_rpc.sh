#!/bin/bash

SVR=localhost:9000
PROTO="./bidrpc/bidrpcproto/product.proto"
SVC="bidrpcproto.ProductService"

# list
echo -e "\n=== ListProducts ==="
grpcurl -plaintext -proto $PROTO -d '{}' $SVR $SVC/ListProducts

# list with pagination
echo -e "\n=== ListProducts (Pagination) ==="
DATA='{"page":1,"page_size":2}'
grpcurl -plaintext -proto $PROTO -d $DATA $SVR $SVC/ListProducts

# 创建产品
echo -e "\n=== CreateProduct ==="
DATA='{"name":"iPhone 15 Pro","description":"Latest iPhone","price":999.99,"quantity":50}'
grpcurl -plaintext -proto $PROTO -d "$DATA" $SVR $SVC/CreateProduct

echo -e "\n=== Get Single Product ==="
DATA='{"id":"6584f023-9cfd-4fe0-b613-2b2ecf00fa4c"}'
grpcurl -plaintext -proto $PROTO -d "$DATA" $SVR $SVC/GetProduct


echo -e "\n=== UpdateProduct  ==="
DATA='{"id":1,"name":"iPhone 15 Pro Max","description":"Updated desc","price":1099.99,"quantity":30}'
grpcurl -plaintext -proto $PROTO -d "$DATA" $SVR $SVC/UpdateProduct


echo -e "\n=== DeleteProduct  ==="
DATA='{"id":"6584f023-9cfd-4fe0-b613-2b2ecf00fa4c"}'
grpcurl -plaintext -proto $PROTO -d "$DATA" $SVR $SVC/DeleteProduct
