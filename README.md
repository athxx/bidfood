# Bidfood grpc Microservices

## Overview

this markdown only for question 1, ratelimter part please view this file
[ratelimiter/README.md](ratelimiter/README.md)

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

## Requirements

- **Go Version**: 1.24.1 +
- **Protoc Version**: 3.15+
- **Dependencies** : Make

## How to run ?

### install dependencies

- 1. install `go`, you can install golang by using `mise`
- 2. install `protoc`
- 3. after install `go` and `protoc`, run `make deps` to generate protobuf files

**please run both api gateway and rpc at the same time, you need to run them in separate terminals.**

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

### unit test

```sh
cd bidone && make test
```

### api test

```sh
cd bidone && make test-api
```

### rpc test

```sh
cd bidone && make test-rpc
```
