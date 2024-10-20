package service

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/router"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/grpc/pagelistServer"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/pb"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func RunFiber(addr string) error {
	app := fiber.New()
	router.Setup(app)

	log.Println("Fiber server is running on port", addr)
	return app.Listen(addr)
}

func RunGrpc(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	pb.RegisterPageListServiceServer(srv, &pagelistServer.Server{})

	log.Println("GRPC server is running on port", addr)
	return srv.Serve(lis)
}

func RunGrpcui(addr, target string) error {
	ctx := context.Background()
	cc, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	h, err := standalone.HandlerViaReflection(ctx, cc, target)
	if err != nil {
		return err
	}

	log.Println("GRPCUI server is running on port", addr)
	return http.ListenAndServe(addr, h)
}
