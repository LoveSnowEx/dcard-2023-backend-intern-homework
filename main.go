package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/router"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/grpc/pagelistServer"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func runFiber(addr string) error {
	app := fiber.New()
	router.Setup(app)

	return app.Listen(addr)
}

func runGrpc(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	pb.RegisterPageListServiceServer(srv, &pagelistServer.Server{})

	return srv.Serve(lis)
}

func initExamplePages() {
	count := 10
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	for i := 1; i <= count; i++ {
		page := page.Page{
			Title:   fmt.Sprintf("Page %d", i),
			Content: fmt.Sprintf("Content %d", i),
			Slug:    fmt.Sprintf("page-%d", i),
		}
		page.ID = uint(i)
		err = dbConn.UpdatePage(&page)
		if err != nil {
			log.Fatalf("failed to create page: %v", err)
		}
	}
}

func main() {
	_, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func() {
		// Close the database connection when the program exits
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}()

	initExamplePages()

	c := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		fmt.Println()
		fmt.Println(sig)
		done <- struct{}{}
	}()
	go func() {
		log.Fatalln(runFiber(":3000"))
		done <- struct{}{}
	}()
	go func() {
		log.Fatalln(runGrpc(":50051"))
		done <- struct{}{}
	}()

	fmt.Println("waiting for Ctrl+C signal")

	<-done
	log.Println("Closing ")
}
