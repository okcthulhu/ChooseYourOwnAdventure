package api

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoryCollection defines the required behavior for interacting with
// the story element-related data in MongoDB. By isolating these methods, we can
// easily swap out the actual MongoDB collection with a mock for testing.
type StoryCollection interface {
	// InsertOne adds a new document to the story elements collection. The method returns
	// the result of the insertion, which contains the ID of the new document, or an error.
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// FindOne locates a single document from the story elements collection based on the filter.
	// A single result is returned, which can be decoded to access the actual document.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{},
		opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// StoryHandler is the main orchestrator for the application's HTTP API. It aggregates
// various dependencies needed to process incoming HTTP requests and produce
// appropriate responses. The fields in this struct adhere to interfaces, thus
// allowing easy substitution for testing and extending functionality.
type StoryHandler struct {
	// StoryCol is an abstraction for the MongoDB collection containing story elements.
	StoryCol StoryCollection
}

// NewStoryHandler serves as a factory function for creating a new instance of the StoryHandler struct.
// It takes in implementations of MongoClientInterface, PlayerCollection, and StoryCollection
// as arguments. By providing these as interfaces, this function allows for greater flexibility
// and testability. For example, you can provide mock implementations when you're writing tests.
// The function returns a pointer to the newly created StoryHandler instance, fully equipped with
// the necessary dependencies for database interactions related to both player and story elements.
func NewStoryHandler(storyCol StoryCollection) *StoryHandler {
	return &StoryHandler{
		StoryCol: storyCol,
	}
}

// CreateStoryElement initializes a new story element in the database with the given details.
// It takes a JSON-formatted request body containing the attributes of the new story element.
// After successful creation, the function returns a JSON-formatted response containing the newly created story element.
// If the operation fails, an appropriate HTTP status code is returned, along with an error message.
func (h *StoryHandler) CreateStoryElement(c echo.Context) error {
	storyElement := new(models.PostStoryElementsJSONRequestBody)
	if err := c.Bind(storyElement); err != nil {
		return err
	}

	// Add this check to handle an empty request body
	if storyElement.IsEmpty() { // Assume you have or will implement an IsEmpty method on your struct
		log.Println("Received empty request body.")
		return c.JSON(http.StatusBadRequest, "Empty request body")
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
func (h *StoryHandler) GetStoryElement(c echo.Context, nodeId string) error {
	filter := bson.M{"nodeID": nodeId}
	singleResult := h.StoryCol.FindOne(context.Background(), filter)
	var storyElement models.StoryElement
	err := singleResult.Decode(&storyElement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, "Story Element not found")
		}
		return c.JSON(http.StatusInternalServerError, "An error occurred")
	}

	return c.JSON(http.StatusOK, storyElement)
}

// UpdateStoryElement modifies an existing story element's state in the database based on the provided updates.
// The function expects a JSON-formatted request body containing the updated attributes of the story element,
// as well as the story element's unique NodeId to identify which record to update.
// Upon successful update, the function returns a JSON-formatted response reflecting the modified story element.
// If the update operation fails or if the specified NodeId does not exist,
// an appropriate HTTP status code and an error message are returned.
func (h *StoryHandler) UpdateStoryElement(c echo.Context, nodeId string, storyElement models.PatchStoryElementsNodeIdJSONRequestBody) error {
	filter := bson.M{"nodeID": nodeId}
	update := bson.M{"$set": storyElement}
	_, err := h.StoryCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Update failed due to an internal error")
	}

	return c.JSON(http.StatusOK, "Story element updated successfully")
}

// DeleteStoryElement removes a story element identified by its node ID from the database.
// It receives an Echo context and the node ID of the story element as parameters.
// The function constructs a filter based on the node ID and attempts to delete the
// corresponding document from the story collection in MongoDB.
// It returns an HTTP status code and a JSON response indicating the outcome of the operation.
// If the deletion is successful, it responds with an HTTP 200 OK status and a success message.
// If an error occurs during the deletion process, it responds with an HTTP 500 Internal Server Error
// status and an error message describing the failure.
func (h *StoryHandler) DeleteStoryElement(c echo.Context, nodeId string) error {
	filter := bson.M{"nodeID": nodeId}
	_, err := h.StoryCol.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Delete failed due to an internal error")
	}

	return c.JSON(http.StatusOK, "Story element deleted successfully")
}
