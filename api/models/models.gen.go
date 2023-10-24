// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package models

// PlayerState defines model for PlayerState.
type PlayerState struct {
	Artifacts      *[]string `json:"artifacts,omitempty"`
	CurrentChapter *string   `json:"currentChapter,omitempty"`
	CurrentPart    *string   `json:"currentPart,omitempty"`
	PlayerId       *string   `json:"playerId,omitempty"`
	Wisdoms        *[]string `json:"wisdoms,omitempty"`
}

// StoryElement defines model for StoryElement.
type StoryElement struct {
	Chapter *struct {
		ArtURL   *string `json:"artURL,omitempty"`
		Name     *string `json:"name,omitempty"`
		VideoURL *string `json:"videoURL,omitempty"`
	} `json:"chapter,omitempty"`
	Content *string                   `json:"content,omitempty"`
	NodeId  *string                   `json:"nodeId,omitempty"`
	Options *[]map[string]interface{} `json:"options,omitempty"`
	Part    *struct {
		ArtURL   *string `json:"artURL,omitempty"`
		Name     *string `json:"name,omitempty"`
		VideoURL *string `json:"videoURL,omitempty"`
	} `json:"part,omitempty"`
	Wisdoms *[]map[string]interface{} `json:"wisdoms,omitempty"`
}
