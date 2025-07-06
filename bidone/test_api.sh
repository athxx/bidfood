#!/bin/bash

ID=6584f023-9cfd-4fe0-b613-2b2ecf00fa4c

# list
echo -e "\n=== ListProducts ==="
curl 'http://localhost:8080/products'

# list with pagination
echo -e "\n=== ListProducts (Pagination) ==="
curl 'http://localhost:8080/products/products?page=1&page_size=1'

# 创建产品
echo -e "\n=== CreateProduct ==="
curl 'http://localhost:8080/products' \
  -X POST \
  -H 'content-type: application/json' \
  -d $'{"name": "iPhone 15 Pro","description": "Latest iPhone with advanced camera system","price": 999.99,"quantity": 50}'

echo -e "\n=== Get Single Product ==="
curl "http://localhost:8080/products/$ID"

echo -e "\n=== UpdateProduct  ==="
curl "http://localhost:8080/products/$ID" \
-X POST \
-H 'content-type: application/json' \
-d '{"price": 100}'

echo -e "\n=== DeleteProduct  ==="
curl -X DELETE "http://localhost:8080/products/$ID" \
