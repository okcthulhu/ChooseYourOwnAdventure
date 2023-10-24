package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api"
)

func main() {
	port := flag.String("port", "1323", "Port to run the server on")
	flag.Parse()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := api.NewHandler()

	e.GET("/adventure", handler.GetAdventure)
	e.PUT("/adventure", handler.UpdateAdventure)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", *port)); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	// ...

	// Disconnect MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := handler.DB.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to close MongoDB client: %v", err)
	}
}

// func ham() {
// 	// Set MongoDB client options
// 	clientOptions := options.Client().ApplyURI("mongodb+srv://n8gallenson:Lg2ke370DwkQe7QO@cyoa01.m1d2ueq.mongodb.net/?retryWrites=true&w=majority")

// 	// Connect to MongoDB
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}

// 	// Check the connection
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to ping MongoDB: %v", err)
// 	} else {
// 		fmt.Println("Successfully connected to Atlas")
// 	}

// 	// Disconnect from MongoDB
// 	defer func() {
// 		if err := client.Disconnect(ctx); err != nil {
// 			log.Fatalf("Failed to close MongoDB client: %v", err)
// 		}
// 	}()
// }
