#!/bin/bash

echo "Generating protobuf files..."

protoc --proto_path=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       bidrpcproto/product.proto

echo "Protobuf generation complete!"
