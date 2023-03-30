package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/config"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/service"
)

func initExamplePages(count int) {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	session := &db.DB{DB: dbConn.DB.Begin()}
	for i := 1; i <= count; i++ {
		page := page.Page{
			Title:   fmt.Sprintf("Page %d", i),
			Content: fmt.Sprintf("Content %d", i),
			Slug:    fmt.Sprintf("page-%d", i),
		}
		page.ID = uint(i)
		err = session.UpdatePage(&page)
		if err != nil {
			session.DB.Rollback()
			log.Fatalf("failed to create page: %v", err)
		}
	}
	err = session.DB.Commit().Error
	if err != nil {
		session.DB.Rollback()
		log.Fatalf("failed to commit: %v", err)
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

	initExamplePages(1000)

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
		log.Fatalln(service.RunFiber(":" + config.FiberPort))
		done <- struct{}{}
	}()
	go func() {
		log.Fatalln(service.RunGrpc(":" + config.GrpcPort))
		done <- struct{}{}
	}()
	go func() {
		log.Fatalln(service.RunGrpcui(":"+config.GrpcuiPort, ":"+config.GrpcPort))
		done <- struct{}{}
	}()

	fmt.Println("waiting for Ctrl+C signal")

	<-done
	log.Println("Closing ")
}
