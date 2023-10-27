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

// MongoClientInterface serves as an abstraction over the native MongoDB client.
// It exposes just enough functionality to perform essential database operations,
// thereby adhering to the Interface Segregation Principle.
type MongoClientInterface interface {
	// Database gets a handle for a MongoDB database with the given name.
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

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

// StoryCollection behaves similarly to PlayerCollection but is intended
// for story elements. This separation makes the system more modular and adheres
// to the Single Responsibility and Interface Segregation Principles.
type StoryCollection interface {
	// InsertOne adds a new document to the story elements collection. The method returns
	// the result of the insertion, which contains the ID of the new document, or an error.
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// FindOne locates a single document from the story elements collection based on the filter.
	// A single result is returned, which can be decoded to access the actual document.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// Handler is the main orchestrator for the application's HTTP API. It aggregates
// various dependencies needed to process incoming HTTP requests and produce
// appropriate responses. The fields in this struct adhere to interfaces, thus
// allowing easy substitution for testing and extending functionality.
type Handler struct {
	// DB is an abstraction over the MongoDB client, allowing the handler to
	// interact with any MongoDB database.
	DB MongoClientInterface

	// PlayerCol is an abstraction for the MongoDB collection containing player data.
	PlayerCol PlayerCollection

	// StoryCol is an abstraction for the MongoDB collection containing story elements.
	StoryCol StoryCollection
}

// NewHandler serves as a factory function for creating a new instance of the Handler struct.
// It takes in implementations of MongoClientInterface, PlayerCollection, and StoryCollection
// as arguments. By providing these as interfaces, this function allows for greater flexibility
// and testability. For example, you can provide mock implementations when you're writing tests.
// The function returns a pointer to the newly created Handler instance, fully equipped with
// the necessary dependencies for database interactions related to both player and story elements.
func NewHandler(client MongoClientInterface, playerCol PlayerCollection, storyCol StoryCollection) *Handler {
	return &Handler{
		DB:        client,
		PlayerCol: playerCol,
		StoryCol:  storyCol,
	}
}

// CreatePlayerState initializes a new player state in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new player state.
// After successful creation, the function returns a JSON-formatted response containing the newly created player state.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *Handler) CreatePlayerState(c echo.Context) error {
	playerState := new(models.PostPlayerJSONRequestBody)
	if err := c.Bind(playerState); err != nil {
		log.Println("Failed to bind playerState:", err)
		return err
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
func (h *Handler) GetPlayerStateByWixID(c echo.Context, wixID string) error {
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
func (h *Handler) UpdatePlayerState(c echo.Context, wixID string, playerState models.PatchPlayerPlayerIdJSONRequestBody) error {
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

// CreateStoryElement initializes a new story element in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new story element.
// After successful creation, the function returns a JSON-formatted response containing the newly created story element.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *Handler) CreateStoryElement(c echo.Context) error {
	storyElement := new(models.PostStoryElementsJSONRequestBody)
	if err := c.Bind(storyElement); err != nil {
		return err
	}

	_, err := h.StoryCol.InsertOne(context.Background(), storyElement)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create story element")
	}

	return c.JSON(http.StatusCreated, storyElement)
}

// GetStoryElement retrieves a specific story element identified by its NodeId from the database.
// The function returns a JSON-formatted response containing the details of the story element.
// If the story element is not found in the database, a 404 status code is returned.
func (h *Handler) GetStoryElement(c echo.Context, nodeId string) error {
	filter := bson.M{"nodeID": nodeId}
	singleResult := h.StoryCol.FindOne(context.Background(), filter)
	var storyElement models.StoryElement
	err := singleResult.Decode(&storyElement)
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
	filter := bson.M{"nodeID": nodeId}
	update := bson.M{"$set": storyElement}
	_, err := h.StoryCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Update failed due to an internal error")
	}

	return c.JSON(http.StatusOK, "Story element updated successfully")
}
