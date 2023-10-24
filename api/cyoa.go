package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	DB *mongo.Client
}

func NewHandler() *Handler {
	clientOptions := options.Client().ApplyURI("mongodb+srv://n8gallenson:Lg2ke370DwkQe7QO@cyoa01.m1d2ueq.mongodb.net/?retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	} else {
		fmt.Println("Successfully connected to Atlas")
	}

	return &Handler{
		DB: client,
	}
}

func (h *Handler) GetAdventure(c echo.Context) error {
	// Your code to get adventure data from MongoDB
	return c.String(http.StatusOK, "Adventure retrieved")
}

func (h *Handler) UpdateAdventure(c echo.Context) error {
	// Your code to update adventure data in MongoDB
	return c.String(http.StatusOK, "Adventure updated")
}
