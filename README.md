# go-grpcmw

`go-grpcmw` provides a package and a protobuf generator for managing easily
grpc interceptors.

The package can be used without the protobuf generator. However, using both
together will allow you to avoid writing redundant code.

## Prerequisites

* Protobuf 3.0.0 or later.

## Installation

```shell
go get -u github.com/MarquisIO/go-grpcmw/protoc-gen-grpc-middleware
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Quick start

Write your gRPC service definition:

```protobuf
syntax = "proto3";

import "github.com/MarquisIO/go-grpcmw/annotations/annotations.proto";

package pb;

option (grpcmw.package_interceptors) = {
  indexes: ["index"]
};

service SomeService {
  option (grpcmw.service_interceptors) = {
    indexes: ["index"]
  };

  rpc SomeMethod (Message) returns (Message) {
    option (grpcmw.method_interceptors) = {
      indexes: ["index"]
    };
  }
}

message Message {
	string msg = 1;
}
```

Generate the stubs:

```shell
protoc --go_out=plugins=grpc:. --grpc-middleware_out=:. path/to/you/file.proto
```

Use the code generated to add your own middlewares:

```go
// Register an interceptor in the registry
registry.GetClientInterceptor("index").
  AddGRPCUnaryInterceptor(SomeUnaryClientInterceptor).
  AddGRPCStreamInterceptor(SomeStreamClientInterceptor)
registry.GetServerInterceptor("index").
  AddGRPCUnaryInterceptor(SomeUnaryServerInterceptor).
  AddGRPCStreamInterceptor(SomeStreamServerInterceptor)

// Client
clientRouter := grpcmw.NewClientRouter()
clientStub := pb.RegisterClientInterceptors(clientRouter)
clientStub.RegisterSomeService().
	SomeMethod().
	AddGRPCInterceptor(clientUnaryMiddleware)
grpc.Dial(address,
	grpc.WithStreamInterceptor(clientRouter.StreamResolver()),
	grpc.WithUnaryInterceptor(clientRouter.UnaryResolver()),
)

// Server
serverRouter := grpcmw.NewServerRouter()
serverStub := pb.RegisterServerInterceptors(serverRouter)
serverStub.RegisterSomeService().
	SomeMethod().
	AddGRPCInterceptor(serverUnaryMiddleware)
grpc.NewServer(
	grpc.UnaryInterceptor(serverRouter.UnaryResolver()),
	grpc.StreamInterceptor(serverRouter.StreamResolver()),
)
```

## Chaining

Four types of interceptors are provided: `ServerUnaryInterceptor`,
`ServerStreamInterceptor`, `ClientUnaryInterceptor` and
`ClientStreamInterceptor` (corresponding to those defined in
[google.golang.org/grpc](https://godoc.org/google.golang.org/grpc)). They allow
you to chain multiple gRPC interceptors of the same type:

```go
// grpcInterceptor1 -> grpcInterceptor2 -> interceptor1 -> grpcInterceptor3
intcp := grpcmw.NewUnaryServerInterceptor(grpcInterceptor1, grpcInterceptor2).
	AddInterceptor(interceptor1).
	AddGRPCInterceptor(grpcInterceptor3)
```

## Routing

This package also provides a routing feature so that interceptors can be bound
either to:
* a protobuf package: all requests to any service that have been declared in
this package will go through the interceptor.
* a gRPC service: all requests to this service will go through the interceptor.
* a gRPC method: all requests to this method will go through the interceptor.

`ServerRouter` and `ClientRouter` provide one more global level in addition to
the three described above.
These implementations are based on the route construction as defined
[in the official gRPC repository](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md).

```go
serverRouter := grpcmw.NewServerRouter()
pkgInterceptor := grpcmw.NewServerInterceptor("pb")
serviceInterceptor := grpcmw.NewServerInterceptor("Service")

// Interceptors that have been added to `serviceInterceptor` will be called
// each time the gRPC service `Service` will be requested. In other words, all
// requests to "/pb.Service/*" will go through these interceptors.
pkgInterceptor.Register(serviceInterceptor)
serverRouter.GetRegister().
  Register(pkgInterceptor)

// In order to use the router, you have to create the server with it.
grpc.NewServer(
	grpc.UnaryInterceptor(serverRouter.UnaryResolver()),
	grpc.StreamInterceptor(serverRouter.StreamResolver()),
)
```

## Registry

The `registry` package provides an interceptor registry for both server and
client side.

```go
registry.GetServerInterceptor("index").
  AddGRPCUnaryInterceptor(SomeUnaryServerInterceptor).
  AddGRPCStreamInterceptor(SomeStreamServerInterceptor)
```

## Protobuf generation

In order to ease the use of the registry and routing features, a protobuf
generator is provided which you can use with the following command:

```shell
protoc --grpc-middleware_out=:. path/to/you/file.proto
```

### Routing

Say we have the following protobuf file:

```protobuf
syntax = "proto3";

package pb;

service SomeService {
  rpc SomeMethod (Message) returns (Message) {}
}

message Message {
	string msg = 1;
}
```

It will create some helpers for adding interceptors to a package, service or
method.

```go
serverRouter := grpcmw.NewServerRouter()
serverStub := pb.RegisterServerInterceptors(serverRouter)
serverStub.AddGRPCInterceptor(pkgUnaryMiddleware)
serviceStub := serverStub.RegisterSomeService()
serviceStub.AddGRPCInterceptor(serviceUnaryMiddleware)
methodStub := serviceStub.SomeMethod()
methodStub.AddGRPCInterceptor(methodUnaryMiddleware)
```

### Registry

Three annotations are provided in
[annotations/annotations.proto](./annotations/annotations.proto):
* `package_interceptors`: for the package level.
* `service_interceptors`: for the service level.
* `method_interceptors`: for the method level.

These annotations have an array of index (`indexes`) that tells the generator
which interceptors from the registry have to be added to the router.

Say we have the following protobuf file:

```protobuf
syntax = "proto3";

import "github.com/MarquisIO/go-grpcmw/annotations/annotations.proto";

package pb;

option (grpcmw.package_interceptors) = {
  indexes: ["index"]
};

service SomeService {
  option (grpcmw.service_interceptors) = {
    indexes: ["index"]
  };

  rpc SomeMethod (Message) returns (Message) {
    option (grpcmw.method_interceptors) = {
      indexes: ["index"]
    };
  }
}

message Message {
	string msg = 1;
}
```

You can then register interceptors in the registry at the index "index".

```go
// Register an interceptor in the registry
registry.GetServerInterceptor("index").
  AddGRPCUnaryInterceptor(SomeUnaryServerInterceptor).
  AddGRPCStreamInterceptor(SomeStreamServerInterceptor)

// You have to call `RegisterServerInterceptors` and `RegisterSomeService` so
// that interceptors are added to the router at the package, servive and method
// levels.
serverRouter := grpcmw.NewServerRouter()
serverStub := pb.RegisterServerInterceptors(serverRouter)
serverStub.RegisterSomeService()
```
