package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize handler with DB connection
	mongoURI := os.Getenv("MONGO_URI") // Read the URI from an environment variable
	if mongoURI == "" {
		log.Fatal("Environment variable MONGO_URI not set")
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetRegistry(api.MongoRegistry)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
	}

	playerCol := client.Database("cyoa").Collection("players")
	storyCol := client.Database("cyoa").Collection("storyElements")
	playerHandler := api.NewPlayerHandler(playerCol)
	storyHandler := api.NewStoryHandler(storyCol)

	// Handle port flag
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the application on")
	flag.Parse()

	// Initialize Echo
	e := echo.New()

	// Define the routes
	// General route comes first

	//Player routes
	e.POST("/player", playerHandler.CreatePlayerState)
	e.GET("/player/:wixID", func(c echo.Context) error {
		wixID := c.Param("wixID")
		return playerHandler.GetPlayerStateByWixID(c, wixID)
	})
	e.PATCH("/player/:wixID", func(c echo.Context) error {
		wixID := c.Param("wixID")
		playerState := new(models.PatchPlayerPlayerIdJSONRequestBody)
		if err := c.Bind(playerState); err != nil {
			return err
		}
		return playerHandler.UpdatePlayerState(c, wixID, *playerState)
	})

	// StoryElement routes
	e.POST("/storyElements", storyHandler.CreateStoryElement)
	e.GET("/storyElements/:nodeId", func(c echo.Context) error {
		return storyHandler.GetStoryElement(c, c.Param("nodeId"))
	})
	e.PUT("/storyElements/:nodeId", func(c echo.Context) error {
		storyElement := new(models.StoryElement)
		if err := c.Bind(storyElement); err != nil {
			return err
		}
		return storyHandler.UpdateStoryElement(c, c.Param("nodeId"), *storyElement)
	})

	// Start the Echo web server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
