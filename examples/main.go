package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/MarquisIO/go-grpcmw/examples/proto"
	"github.com/MarquisIO/go-grpcmw/examples/server"
	"github.com/MarquisIO/go-grpcmw/grpcmw"
	"github.com/MarquisIO/go-grpcmw/grpcmw/registry"
	"google.golang.org/grpc"
)

func serverMiddlewareRegistry(level string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("enter server : %s level of middleware (registry)\n", level)
		defer fmt.Printf("leave server : %s level of middleware (registry)\n", level)
		return handler(ctx, req)
	}
}

func clientMiddlewareRegistry(level string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Printf("enter client : %s level of middleware (registry)\n", level)
		defer fmt.Printf("leave client : %s level of middleware (registry)\n", level)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func serverMiddleware(level string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("enter server : %s level of middleware\n", level)
		defer fmt.Printf("leave server : %s level of middleware\n", level)
		return handler(ctx, req)
	}
}

func clientMiddleware(level string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fmt.Printf("enter client : %s level of middleware\n", level)
		defer fmt.Printf("leave client : %s level of middleware\n", level)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func startServer(port uint16) (*grpc.Server, net.Listener) {
	// Server
	// Setup global server router
	r := grpcmw.NewServerRouter()
	r.GetRegister().AddGRPCUnaryInterceptor(serverMiddleware("global"))

	pkgInterceptors := pb.RegisterServerInterceptors(r)
	pkgInterceptors.AddGRPCUnaryInterceptor(serverMiddleware("package"))
	pkgInterceptors.RegisterService().AddGRPCUnaryInterceptor(serverMiddleware("service"))
	pkgInterceptors.RegisterService().Method().AddGRPCInterceptor(serverMiddleware("method"))

	// Create gRPC server and register the service
	var e serverpb.Example
	server := grpc.NewServer(grpc.UnaryInterceptor(r.UnaryResolver()))
	pb.RegisterServiceServer(server, &e)

	// Start listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Could not create listener on port %d: %v", port, err)
	}
	go server.Serve(lis)
	return server, lis
}

func startClient(port uint16) (*grpc.ClientConn, pb.ServiceClient) {
	// Client
	// Setup global client router
	r := grpcmw.NewClientRouter()
	r.GetRegister().AddGRPCUnaryInterceptor(clientMiddleware("global"))

	pkgInterceptors := pb.RegisterClientInterceptors(r)
	pkgInterceptors.AddGRPCUnaryInterceptor(clientMiddleware("package"))
	pkgInterceptors.RegisterService().AddGRPCUnaryInterceptor(clientMiddleware("service"))
	pkgInterceptors.RegisterService().Method().AddGRPCInterceptor(clientMiddleware("method"))

	// Setup connection to the server
	target := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(r.UnaryResolver()))
	if err != nil {
		log.Fatalf("Could not dial \"%s\": %v", target, err)
	}
	return conn, pb.NewServiceClient(conn)
}

func main() {
	var port uint16 = 4242

	// Register middlewares on registry
	registry.GetClientInterceptor("pkg").AddGRPCUnaryInterceptor(clientMiddlewareRegistry("package"))
	registry.GetClientInterceptor("srv").AddGRPCUnaryInterceptor(clientMiddlewareRegistry("service"))
	registry.GetClientInterceptor("meth").AddGRPCUnaryInterceptor(clientMiddlewareRegistry("method"))
	registry.GetServerInterceptor("pkg").AddGRPCUnaryInterceptor(serverMiddlewareRegistry("package"))
	registry.GetServerInterceptor("srv").AddGRPCUnaryInterceptor(serverMiddlewareRegistry("service"))
	registry.GetServerInterceptor("meth").AddGRPCUnaryInterceptor(serverMiddlewareRegistry("method"))

	server, lis := startServer(port)
	defer lis.Close()
	defer server.GracefulStop()

	conn, client := startClient(port)
	defer conn.Close()

	msg, err := client.Method(context.Background(), &pb.Message{Msg: "message"})
	if err != nil {
		log.Fatalf("Call to Method failed: %v", err)
	}
	fmt.Printf("Received : %s\n", msg.Msg)
}
