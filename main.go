package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/config"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"
)

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b)
}

func initExamplePages(count int64) {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	var existCount int64
	dbConn.DB.Model(&page.Page{}).Count(&existCount)

	pages := make([]*page.Page, 0, count)
	for i := existCount; i < int64(count); i++ {
		page := page.Page{
			Title:   fmt.Sprintf("Page %d", i),
			Content: fmt.Sprintf("Content %d", i),
			Slug:    fmt.Sprintf("page-%s", randomString(8)),
		}
		pages = append(pages, &page)
	}
	dbConn.DB.CreateInBatches(&pages, 1000)
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

	conf := config.Get()

	_ = conf

	go func() {
		sig := <-c
		fmt.Println()
		fmt.Println(sig)
		done <- struct{}{}
	}()
	// go func() {
	// 	log.Fatalln(service.RunFiber(":" + conf.FiberPort))
	// 	done <- struct{}{}
	// }()
	// go func() {
	// 	log.Fatalln(service.RunGrpc(":" + conf.GrpcPort))
	// 	done <- struct{}{}
	// }()
	// go func() {
	// 	log.Fatalln(service.RunGrpcui(":"+conf.GrpcuiPort, ":"+conf.GrpcPort))
	// 	done <- struct{}{}
	// }()

	fmt.Println("waiting for Ctrl+C signal")

	<-done
	log.Println("Closing ")
}
