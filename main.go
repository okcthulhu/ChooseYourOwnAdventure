package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
)

func main() {
	// Handle port flag
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the application on")
	flag.Parse()

	// Initialize Echo
	e := echo.New()

	// Initialize handler with DB connection
	handler := api.NewHandler()

	// Define the routes

	// General route comes first

	//Player routes
	// e.POST("/player", handler.CreatePlayerState)
	// e.GET("/player/{playerId}", func(c echo.Context) error {
	// 	return handler.GetPlayerStateByWixID(c, c.Param("playerId"))
	// })
	// e.PATCH("/player/:playerId", func(c echo.Context) error {
	// 	playerState := new(models.Player)
	// 	if err := c.Bind(playerState); err != nil {
	// 		return err
	// 	}
	// 	return handler.UpdatePlayerState(c, c.Param("playerId"), *playerState)
	// })
	e.POST("/player", handler.CreatePlayerState)
	e.GET("/player/:wixID", func(c echo.Context) error {
		wixID := c.Param("wixID")
		return handler.GetPlayerStateByWixID(c, wixID)
	})
	e.PATCH("/player/:wixID", func(c echo.Context) error {
		wixID := c.Param("wixID")
		playerState := new(models.PatchPlayerPlayerIdJSONRequestBody)
		if err := c.Bind(playerState); err != nil {
			return err
		}
		return handler.UpdatePlayerState(c, wixID, *playerState)
	})

	// StoryElement routes
	e.POST("/storyElements", handler.CreateStoryElement)
	e.GET("/storyElements/:nodeId", func(c echo.Context) error {
		return handler.GetStoryElement(c, c.Param("nodeId"))
	})
	e.PUT("/storyElements/:nodeId", func(c echo.Context) error {
		storyElement := new(models.StoryElement)
		if err := c.Bind(storyElement); err != nil {
			return err
		}
		return handler.UpdateStoryElement(c, c.Param("nodeId"), *storyElement)
	})

	// Start the Echo web server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
