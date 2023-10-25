package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	DB *mongo.Client
}

// NewHandler initializes a new Handler struct with a MongoDB client.
// It establishes a connection to the MongoDB database by applying the given URI.
// A context with a timeout is also set up to handle the database connection.
// The function returns a pointer to the newly created Handler, which includes the MongoDB client.
// If the connection to MongoDB fails, an error message is printed to the console.
func NewHandler() *Handler {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://n8gallenson:Lg2ke370DwkQe7QO@cyoa01.m1d2ueq.mongodb.net/?retryWrites=true&w=majority").
		SetRegistry(mongoRegistry)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
	}

	return &Handler{
		DB: client,
	}
}

// CreatePlayerState initializes a new player state in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new player state.
// After successful creation, the function returns a JSON-formatted response containing the newly created player state.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *Handler) CreatePlayerState(c echo.Context) error {
	playerState := new(models.PostPlayerJSONRequestBody)
	if err := c.Bind(playerState); err != nil {
		return err
	}

	collection := h.DB.Database("cyoa").Collection("player")

	_, err := collection.InsertOne(context.Background(), playerState)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, "Failed to create player state")
	}

	return c.JSON(http.StatusCreated, playerState)
}

// GetPlayerState fetches the current state of a player from the database
// given a playerId. The function returns a JSON response containing the
// player's state if found, or a 404 status code if the player is not found.
func (h *Handler) GetPlayerState(c echo.Context, playerId string) error {
	collection := h.DB.Database("cyoa").Collection("player")
	var playerState models.Player

	// Convert string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(playerId)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	filter := bson.M{"_id": objectId}

	err = collection.FindOne(context.Background(), filter).Decode(&playerState)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusNotFound, "Player not found")
	}

	return c.JSON(http.StatusOK, playerState)
}

// GetPlayerStateByWixID fetches the current state of a player from the database
// given a WixID. The function returns a JSON response containing the
// player's state if found, or a 404 status code if the player is not found.
func (h *Handler) GetPlayerStateByWixID(c echo.Context, wixID string) error {
	collection := h.DB.Database("cyoa").Collection("player")
	var playerState models.Player
	filter := bson.M{"wixID": wixID}

	err := collection.FindOne(context.Background(), filter).Decode(&playerState)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusNotFound, "Player not found")
	}

	return c.JSON(http.StatusOK, playerState)
}

// GetPlayerStateByUsername fetches the current state of a player from the database
// given a username. The function returns a JSON response containing the
// player's state if found, or a 404 status code if the player is not found.
func (h *Handler) GetPlayerStateByUsername(c echo.Context) error {
	username := c.QueryParam("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, "Username is required")
	}

	collection := h.DB.Database("cyoa").Collection("player")
	var playerState models.Player
	filter := bson.M{"username": username}

	err := collection.FindOne(context.Background(), filter).Decode(&playerState)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(http.StatusNotFound, "Player not found")
	}

	return c.JSON(http.StatusOK, playerState)
}

// UpdatePlayerState modifies an existing player's state in the database based on the provided updates.
// The function expects a JSON-formatted request body containing the updated attributes of the player state,
// as well as the player's unique ID to identify which record to update.
// Upon successful update, the function returns a JSON-formatted response reflecting the modified player state.
// If the update operation fails or if the specified player ID does not exist,
// an appropriate HTTP status code and an error message are returned.
func (h *Handler) UpdatePlayerState(c echo.Context, playerId string, playerState models.PatchPlayerPlayerIdJSONRequestBody) error {
	collection := h.DB.Database("cyoa").Collection("player")

	// Convert string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(playerId)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": playerState}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusNotFound, "Player not found or update failed")
	}

	return c.JSON(http.StatusOK, "Player updated successfully")
}

// CreateStoryElement initializes a new story element in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new story element.
// After successful creation, the function returns a JSON-formatted response containing the newly created story element.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *Handler) CreateStoryElement(c echo.Context) error {
	storyElement := new(models.PostStoryElementsJSONRequestBody)
	if err := c.Bind(storyElement); err != nil {
		return err
	}

	collection := h.DB.Database("cyoa").Collection("storyElements")

	_, err := collection.InsertOne(context.Background(), storyElement)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, "Failed to create story element")
	}

	return c.JSON(http.StatusCreated, storyElement)
}

// GetStoryElement retrieves a specific story element identified by its NodeId from the database.
// The function returns a JSON-formatted response containing the details of the story element.
// If the story element is not found in the database, a 404 status code is returned.
func (h *Handler) GetStoryElement(c echo.Context, nodeId string) error {
	collection := h.DB.Database("cyoa").Collection("storyElements")
	var storyElement models.StoryElement
	filter := bson.M{"nodeID": nodeId}

	err := collection.FindOne(context.Background(), filter).Decode(&storyElement)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Story Element not found")
	}

	return c.JSON(http.StatusOK, storyElement)
}

// UpdateStoryElement modifies an existing story element's state in the database based on the provided updates.
// The function expects a JSON-formatted request body containing the updated attributes of the story element,
// as well as the story element's unique NodeId to identify which record to update.
// Upon successful update, the function returns a JSON-formatted response reflecting the modified story element.
// If the update operation fails or if the specified NodeId does not exist,
// an appropriate HTTP status code and an error message are returned.
func (h *Handler) UpdateStoryElement(c echo.Context, nodeId string, storyElement models.PutStoryElementsNodeIdJSONRequestBody) error {
	collection := h.DB.Database("cyoa").Collection("storyElements")

	filter := bson.M{"nodeID": nodeId}
	update := bson.M{"$set": storyElement}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusNotFound, "Story element not found or update failed")
	}

	return c.JSON(http.StatusOK, "Story element updated successfully")
}
