package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PlayerCollection defines the required behavior for interacting with
// the player-related data in MongoDB. By isolating these methods, we can
// easily swap out the actual MongoDB collection with a mock for testing.
type PlayerCollection interface {
	// InsertOne adds a new document to the players collection. It returns the
	// result of the insertion operation, which includes the ID of the newly
	// inserted document, or an error if the operation fails.
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// FindOne searches for a single document in the players collection that matches
	// the filter. The method returns a single result which can be decoded to
	// obtain the document's data.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// PlayerHandler is the main orchestrator for the application's HTTP API. It aggregates
// various dependencies needed to process incoming HTTP requests and produce
// appropriate responses. The fields in this struct adhere to interfaces, thus
// allowing easy substitution for testing and extending functionality.
type PlayerHandler struct {
	// PlayerCol is an abstraction for the MongoDB collection containing player data.
	PlayerCol PlayerCollection
}

// NewPlayerHandler serves as a factory function for creating a new instance of the PlayerHandler struct.
// It takes in implementations of MongoClientInterface, PlayerCollection, and StoryCollection
// as arguments. By providing these as interfaces, this function allows for greater flexibility
// and testability. For example, you can provide mock implementations when you're writing tests.
// The function returns a pointer to the newly created PlayerHandler instance, fully equipped with
// the necessary dependencies for database interactions related to both player and story elements.
func NewPlayerHandler(playerCol PlayerCollection) *PlayerHandler {
	return &PlayerHandler{
		PlayerCol: playerCol,
	}
}

// CreatePlayerState initializes a new player state in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new player state.
// After successful creation, the function returns a JSON-formatted response containing the newly created player state.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *PlayerHandler) CreatePlayerState(c echo.Context) error {
	playerState := new(models.PostPlayerJSONRequestBody)
	if err := c.Bind(playerState); err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to bind the request to the player")
	}

	// Add this check to handle an empty request body
	if playerState.IsEmpty() { // Assume you have or will implement an IsEmpty method on your struct
		return c.JSON(http.StatusBadRequest, "Empty request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.PlayerCol.InsertOne(ctx, playerState)
	if err != nil {
		log.Println("Failed to insert player state:", err)
		return c.JSON(http.StatusInternalServerError, "Failed to create player state")
	}

	return c.JSON(http.StatusCreated, playerState)
}

// GetPlayerStateByWixID fetches the current state of a player from the database
// given a WixID. The function returns a JSON response containing the
// player's state if found, or a 404 status code if the player is not found.
func (h *PlayerHandler) GetPlayerStateByWixID(c echo.Context, wixID string) error {
	parsedUUID, err := uuid.Parse(wixID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid WixID format")
	}

	binaryUUID := primitive.Binary{
		Subtype: 0x04,
		Data:    parsedUUID[:],
	}

	filter := bson.M{"wixID": binaryUUID}
	singleResult := h.PlayerCol.FindOne(context.Background(), filter)
	var playerState models.Player
	err = singleResult.Decode(&playerState)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Player not found")
	}

	return c.JSON(http.StatusOK, playerState)
}

// UpdatePlayerState modifies an existing player's state in the database based on the provided updates.
// The function expects a JSON-formatted request body containing the updated attributes of the player state,
// as well as the player's Wix ID to identify which record to update.
// Upon successful update, the function returns a JSON-formatted response reflecting the modified player state.
// If the update operation fails or if the specified Wix ID does not exist,
// an appropriate HTTP status code and an error message are returned.
func (h *PlayerHandler) UpdatePlayerState(c echo.Context, wixID string, playerState models.PatchPlayerPlayerIdJSONRequestBody) error {
	parsedUUID, err := uuid.Parse(wixID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid WixID format")
	}

	binaryUUID := primitive.Binary{
		Subtype: 0x04,
		Data:    parsedUUID[:],
	}

	filter := bson.M{"wixID": binaryUUID}
	update := bson.M{"$set": playerState}
	_, err = h.PlayerCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Player not found or update failed")
	}

	return c.JSON(http.StatusOK, "Player updated successfully")
}
