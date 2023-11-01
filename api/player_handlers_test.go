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
		artifacts := []string{"artifact1", "artifact2"}
		currentChapter := "chapter1"
		currentPart := "part1"
		storyID := "someStoryId"
		wisdoms := []string{"wisdom1", "wisdom2"}

		storyState := models.StoryState{
			Artifacts:      &artifacts,
			CurrentChapter: &currentChapter,
			CurrentPart:    &currentPart,
			StoryID:        &storyID,
			Wisdoms:        &wisdoms,
		}

		storyStates := []models.StoryState{storyState}

		email := "test@example.com"
		username := "TestUser"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayerJSONRequestBody{
			Email:       &email,
			StoryStates: &storyStates,
			Username:    &username,
			WixID:       &wixID,
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
		email := "test@example.com"
		username := "TestUser"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayerJSONRequestBody{
			Email:    &email,
			Username: &username,
			WixID:    &wixID,
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
		email := "test@example.com"
		username := "TestUser"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PostPlayerJSONRequestBody{
			Email:    &email,
			Username: &username,
			WixID:    &wixID,
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
		expectedResponse := fmt.Sprintf("{\"email\":\"%s\",\"wixID\":\"%s\"}", *playerState.Email, playerState.WixID.String())
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

func TestUpdatePlayerState_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player state updated", func(mt *mtest.T) {
		email := "test@example.com"
		username := "TestUser"
		wixID := uuid.New()

		// Create request payload
		playerState := &models.PatchPlayerPlayerIdJSONRequestBody{
			Email:    &email,
			Username: &username,
			WixID:    &wixID,
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

func TestUpdatePlayerState_InvalidWixID(t *testing.T) {
	wixID := "invalidWixID"

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/player/%s", wixID), nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("wixID", wixID)

	h := api.NewPlayerHandler(nil) // Passing nil as we expect this to fail before DB interaction

	h.UpdatePlayerState(c, wixID, models.PatchPlayerPlayerIdJSONRequestBody{})
	assert.Equal(t, "\"Invalid WixID format\"", strings.TrimSuffix(rec.Body.String(), "\n"))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePlayerState_FailedUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("player update failed", func(mt *mtest.T) {
		wixID := uuid.New()
		playerState := models.PatchPlayerPlayerIdJSONRequestBody{
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
