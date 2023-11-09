package api_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// CreateStoryElement
func TestCreateStoryElement_StoryElementCreated(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("story element created", func(mt *mtest.T) {
		// Test data
		storyID := "Sample Story ID"
		content := "This is sample content."

		// Create request payload
		storyElement := &models.PostStoryElementsJSONRequestBody{
			StoryID: storyID,
			Content: content,
		}

		storyElementJSON, err := json.Marshal(storyElement)
		if err != nil {
			log.Fatalf("Failed to serialize storyElement: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/story", bytes.NewBuffer(storyElementJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		h := api.NewStoryHandler(mt.Coll)

		// Add Mock Response
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})

		h.CreateStoryElement(c)

		// Validate
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestCreateStoryElement_EmptyRequestBody(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("empty request body", func(mt *mtest.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/story", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := api.NewStoryHandler(mt.Coll)
		h.CreateStoryElement(c)

		// Validate
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `"Empty request body"`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestCreateStoryElement_InsertFailed(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("insertion failed", func(mt *mtest.T) {
		// Test data
		storyID := "Failed Story ID"
		content := "Failed content."

		// Create request payload
		storyElement := &models.PostStoryElementsJSONRequestBody{
			StoryID: storyID,
			Content: content,
		}

		storyElementJSON, err := json.Marshal(storyElement)
		if err != nil {
			log.Fatalf("Failed to serialize storyElement: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/story", bytes.NewBuffer(storyElementJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		h := api.NewStoryHandler(mt.Coll)

		// Add Mock Response for insertion failure
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "insertion error"}})

		h.CreateStoryElement(c)

		// Validate
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, `"Failed to create story element"`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

// GetStoryElement

func TestGetStoryElement_StoryElementFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("story element found", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		req := httptest.NewRequest(http.MethodGet, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		content := "Some content"
		storyElement := &models.StoryElement{
			NodeID:  nodeId,
			Content: content,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "nodeID", Value: storyElement.NodeID},
			{Key: "content", Value: storyElement.Content},
		}))

		h.GetStoryElement(c, nodeId)

		// Validate
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"nodeID":"SomeNodeID"`)
		assert.Contains(t, rec.Body.String(), `"content":"Some content"`)
	})
}

func TestGetStoryElement_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("story element not found", func(mt *mtest.T) {
		nodeId := "NonExistentNodeID"
		req := httptest.NewRequest(http.MethodGet, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		h.GetStoryElement(c, nodeId)

		// Validate
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Story Element not found")
	})
}

func TestGetStoryElement_InternalServerError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("internal server error", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		req := httptest.NewRequest(http.MethodGet, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "Internal Server Error"}})

		h.GetStoryElement(c, nodeId)

		// Validate
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "An error occurred")
	})
}

func TestGetStoryElement_InvalidNodeID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("invalid nodeID", func(mt *mtest.T) {
		nodeId := ""
		req := httptest.NewRequest(http.MethodGet, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		// Assuming that the MongoDB response would be empty for invalid nodeID
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		h.GetStoryElement(c, nodeId)

		// Validate
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Story Element not found")
	})
}

// UpdateStoryElement

func TestUpdateStoryElement_SuccessfulUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("successful update", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		content := "New Content"
		storyElement := models.PatchStoryElementsNodeIdJSONRequestBody{Content: content, NodeID: nodeId}
		reqBody, _ := json.Marshal(storyElement)
		req := httptest.NewRequest(http.MethodPut, "/story/"+nodeId, bytes.NewBuffer(reqBody))
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		h.UpdateStoryElement(c, nodeId, storyElement)

		// Validate
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Story element updated successfully")
	})
}

func TestUpdateStoryElement_InternalServerError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("internal server error", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		updatedContent := "Updated content"
		storyElement := models.PatchStoryElementsNodeIdJSONRequestBody{Content: updatedContent}

		req := httptest.NewRequest(http.MethodPut, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("nodeId", nodeId)
		c.Set("storyElement", storyElement)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "Internal Server Error"}})

		h.UpdateStoryElement(c, nodeId, storyElement)

		// Validate
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Update failed due to an internal error")
	})
}

// DeleteStoryElement

func TestDeleteStoryElement_Deleted(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("story element deleted successfully", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		req := httptest.NewRequest(http.MethodDelete, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/story/:nodeId")
		c.Set("nodeId", nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := h.DeleteStoryElement(c, nodeId)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Story element deleted successfully")
	})
}

func TestDeleteStoryElement_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("story element not found", func(mt *mtest.T) {
		nodeId := "NonExistentNodeID"
		req := httptest.NewRequest(http.MethodDelete, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/story/:nodeId")
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(bson.D{{Key: "n", Value: 0}, {Key: "ok", Value: 1}})

		err := h.DeleteStoryElement(c, nodeId)
		// Since the MongoDB driver does not return an error for delete operations
		// when no document is found, we don't expect an error here.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		// Optionally, check for some indication of a no-op
	})
}

func TestDeleteStoryElement_InternalServerError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("internal server error", func(mt *mtest.T) {
		nodeId := "SomeNodeID"
		req := httptest.NewRequest(http.MethodDelete, "/story/"+nodeId, nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("/story/:nodeId")
		c.SetParamNames("nodeId")
		c.SetParamValues(nodeId)

		h := api.NewStoryHandler(mt.Coll)

		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "Internal Server Error"}})

		h.DeleteStoryElement(c, nodeId)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Delete failed due to an internal error")
	})
}
