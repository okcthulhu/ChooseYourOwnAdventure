package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api"
	"github.com/okcthulhu/ChooseYourOwnAdventure/api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// CreatePlayerState

func TestCreatePlayerState_PlayerCreated(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player created", func(mt *mtest.T) {
		currentStoryNodeID := "some story node ID"
		storyID := "someStoryId"
		wisdoms := []models.Wisdom{{
			Name:     "wisdom 1",
			WisdomID: "wisdom id 1",
		}}

		storyState := models.StoryState{
			CurrentStoryNodeID: currentStoryNodeID,
			StoryID:            storyID,
			Wisdoms:            &wisdoms,
		}

		storyStates := []models.StoryState{storyState}

		var email openapi_types.Email = "test@example.com"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayersJSONRequestBody{
			Email:       email,
			StoryStates: &storyStates,
			WixID:       wixID,
		}

		playerStateJSON, err := json.Marshal(playerState)
		if err != nil {
			log.Fatalf("Failed to serialize playerState: %v", err)
		}

		// Create request using Echo's methods
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(playerStateJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Create a Response recorder
		rec := httptest.NewRecorder()

		// Create Echo context
		e := echo.New()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"wixID", playerState.WixID},
			{"email", playerState.Email},
		}))

		h.CreatePlayerState(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestCreatePlayerState_InsertFailed(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("failed insertion into MongoDB", func(mt *mtest.T) {
		var email openapi_types.Email = "test@example.com"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayersJSONRequestBody{
			Email: email,
			WixID: wixID,
		}

		playerStateJSON, err := json.Marshal(playerState)
		if err != nil {
			log.Fatalf("Failed to serialize playerState: %v", err)
		}

		// Create request using Echo's methods
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(playerStateJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Create a Response recorder
		rec := httptest.NewRecorder()

		// Create Echo context
		e := echo.New()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)

		// Simulate an insert failure
		mt.AddMockResponses(bson.D{{"ok", 0}, {"errmsg", "insertion error"}})

		h.CreatePlayerState(c)

		assert.Equal(t, `"Failed to create player state"`, strings.TrimSuffix(rec.Body.String(), "\n"))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestCreatePlayerState_EmptyRequestBody(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("empty request body", func(mt *mtest.T) {
		e := echo.New()
		req := httptest.NewRequest("POST", "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)

		h.CreatePlayerState(c)

		// Check the response code
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Check the response body
		assert.Equal(t, `"Empty request body"`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestCreatePlayerState_FieldTypeMismatch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("fail to bind request body", func(mt *mtest.T) {
		// Create a malformed JSON request
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("malformed json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Create a Response recorder
		rec := httptest.NewRecorder()

		// Create Echo context
		e := echo.New()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)

		h.CreatePlayerState(c)
		assert.Equal(t, `"Failed to bind the request to the player"`, strings.TrimSuffix(rec.Body.String(), "\n"))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// GetPlayerStateByWixID

func TestGetPlayerStateByWixID_PlayerFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player found", func(mt *mtest.T) {
		var email openapi_types.Email = "test@example.com"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayersJSONRequestBody{
			Email: email,
			WixID: wixID,
		}

		req := httptest.NewRequest("GET", fmt.Sprintf("/player/%s", wixID), nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("wixID", wixID.String())

		h := api.NewPlayerHandler(mt.Coll)

		// Use a string as the mock return value for "wixID" instead of binary.
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"wixID", playerState.WixID},
			{"email", playerState.Email},
		}))

		h.GetPlayerStateByWixID(c, wixID.String())

		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := fmt.Sprintf("{\"email\":\"%s\",\"wixID\":\"%s\"}", playerState.Email, playerState.WixID.String())
		assert.Equal(t, expectedResponse, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestGetPlayerStateByWixID_InvalidWixID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("invalid wixID", func(mt *mtest.T) {
		wixID := "invalidWixID"

		req := httptest.NewRequest("GET", fmt.Sprintf("/player/%s", wixID), nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("wixID", wixID)

		h := api.NewPlayerHandler(mt.Coll)
		h.GetPlayerStateByWixID(c, wixID)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "\"Invalid WixID format\"", strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestGetPlayerStateByWixID_PlayerNotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player not found", func(mt *mtest.T) {
		wixID := uuid.New()

		req := httptest.NewRequest("GET", fmt.Sprintf("/player/%s", wixID), nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("wixID", wixID.String())

		h := api.NewPlayerHandler(mt.Coll)
		h.GetPlayerStateByWixID(c, wixID.String())

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "\"Player not found\"", strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

// UpdatePlayerState

func TestUpdatePlayerState_InvalidWixID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("invalid WixID format", func(mt *mtest.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)
		h.UpdatePlayerState(c, "invalidUUID", models.PatchPlayersPlayerIdJSONRequestBody{})

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid WixID format")
	})
}

func TestUpdatePlayerState_NoStoryStates(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("no story states provided", func(mt *mtest.T) {
		e := echo.New()
		wixID := uuid.New().String()
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := api.NewPlayerHandler(mt.Coll)
		h.UpdatePlayerState(c, wixID, models.PatchPlayersPlayerIdJSONRequestBody{})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "No story states provided")
	})
}

func TestUpdatePlayerState_ValidUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("valid update", func(mt *mtest.T) {
		e := echo.New()
		playerID := uuid.New().String()
		playerWixID := uuid.New()  // Generating a new UUID for the WixID
		storyID := "someStoryID"   // Use a test story ID
		wisdomID := "someWisdomID" // Use a test wisdom ID

		playerUpdate := models.PatchPlayersPlayerIdJSONRequestBody{
			Email: "test@email.com",
			StoryStates: &[]models.StoryState{{
				CurrentStoryNodeID: "testStoryNodeID",
				StoryID:            storyID,
				Wisdoms: &[]models.Wisdom{{
					Name:        "Test Wisdom",
					WisdomID:    wisdomID,
					Description: new(string),
					ArtURL:      new(string),
				}},
			}},
			WixID: playerWixID,
		}

		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock the first update response for setting wisdom details
		// mt.AddMockResponses(mtest.CreateWriteErrorsResponse(errors.New("duplicate"), 11000))

		// Mock the second update response for adding new wisdom
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		h := api.NewPlayerHandler(mt.Coll)

		err := h.UpdatePlayerState(c, playerID, playerUpdate)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Player state updated successfully")
	})
}

func TestUpdatePlayerState_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player state updated", func(mt *mtest.T) {
		var email openapi_types.Email = "test@example.com"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PatchPlayersPlayerIdJSONRequestBody{
			Email: email,
			WixID: wixID,
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/player/%s", wixID), nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("wixID", wixID.String())

		h := api.NewPlayerHandler(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"wixID", playerState.WixID},
			{"email", playerState.Email},
		}))

		err := h.UpdatePlayerState(c, wixID.String(), *playerState)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// func TestUpdatePlayerState_InvalidWixID(t *testing.T) {
// 	wixID := "invalidWixID"

// 	req := httptest.NewRequest("PATCH", fmt.Sprintf("/player/%s", wixID), nil)
// 	rec := httptest.NewRecorder()
// 	e := echo.New()
// 	c := e.NewContext(req, rec)
// 	c.Set("wixID", wixID)

// 	h := api.NewPlayerHandler(nil) // Passing nil as we expect this to fail before DB interaction

// 	h.UpdatePlayerState(c, wixID, models.PatchPlayersPlayerIdJSONRequestBody{})
// 	assert.Equal(t, "\"Invalid WixID format\"", strings.TrimSuffix(rec.Body.String(), "\n"))
// 	assert.Equal(t, http.StatusBadRequest, rec.Code)
// }

func TestUpdatePlayerState_FailedUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player update failed", func(mt *mtest.T) {
		wixID := uuid.New()
		playerState := models.PatchPlayersPlayerIdJSONRequestBody{
			// populate the fields
		}

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/player/%s", wixID), nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.Set("wixID", wixID.String())

		h := api.NewPlayerHandler(mt.Coll)

		mt.AddMockResponses( /* mock to simulate failed UpdateOne */ )

		h.UpdatePlayerState(c, wixID.String(), playerState)
		assert.Equal(t, "\"Player not found or update failed\"", strings.TrimSuffix(rec.Body.String(), "\n"))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
