# Bidfood grpc Microservices

## Overview

this markdown only for question 1, ratelimter part please enter the ratelimiter directory and read the README.md file

### Key Features

- **Microservices Architecture** - Separate gRPC core service and HTTP API gateway
- **Complete CRUD Operations** - Create, read, update, delete products
- **Advanced Querying** - Pagination and name-based filtering
- **Thread-Safe Storage** - In-memory storage with proper synchronization
- **Graceful Shutdown** - Context-based cleanup and resource management
- **Comprehensive Testing** - Unit tests with mocks and integration tests
- **Developer-Friendly** - Rich Makefile with development commands

## Architecture

```
┌─────────────────┐    HTTP     ┌─────────────────┐    gRPC     ┌─────────────────┐
│   Client/Web    │ ◄────────── │   bidapi        │ ◄────────── │   bidrpc        │
│   Application   │    REST     │  (HTTP Gateway) │   Internal  │ (Core Service)  │
└─────────────────┘             └─────────────────┘             └─────────────────┘
                                      :8080                           :9000
```

## How to run ?

### install dependencies

- install `go`, you can install golang by using `mise`
- install `protoc`
- after install `go` and `protoc`, run `make deps` to generate protobuf files

### Run api gateway

```bash

cd bidone && make api
```

### Run rpc service

```bash
cd bidone
# generate .go files from protobuf file
make proto
# run rpc server
make rpc
```

## How to test ?

```sh
cd bidone && make test
```
