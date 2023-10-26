// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new player state.
	// (POST /player)
	PostPlayer(ctx echo.Context) error
	// Retrieve a player's state.
	// (GET /player/{playerId})
	GetPlayerPlayerId(ctx echo.Context, playerId openapi_types.UUID) error
	// Update a player's state by ID.
	// (PATCH /player/{playerId})
	PatchPlayerPlayerId(ctx echo.Context, playerId openapi_types.UUID) error
	// Create a new story element.
	// (POST /storyElements)
	PostStoryElements(ctx echo.Context) error
	// Retrieve a story node.
	// (GET /storyElements/{nodeId})
	GetStoryElementsNodeId(ctx echo.Context, nodeId string) error
	// Update a story node by its ID.
	// (PUT /storyElements/{nodeId})
	PutStoryElementsNodeId(ctx echo.Context, nodeId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostPlayer converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayer(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayer(ctx)
	return err
}

// GetPlayerPlayerId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPlayerPlayerId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "playerId" -------------
	var playerId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "playerId", runtime.ParamLocationPath, ctx.Param("playerId"), &playerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter playerId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPlayerPlayerId(ctx, playerId)
	return err
}

// PatchPlayerPlayerId converts echo context to params.
func (w *ServerInterfaceWrapper) PatchPlayerPlayerId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "playerId" -------------
	var playerId openapi_types.UUID

	err = runtime.BindStyledParameterWithLocation("simple", false, "playerId", runtime.ParamLocationPath, ctx.Param("playerId"), &playerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter playerId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PatchPlayerPlayerId(ctx, playerId)
	return err
}

// PostStoryElements converts echo context to params.
func (w *ServerInterfaceWrapper) PostStoryElements(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostStoryElements(ctx)
	return err
}

// GetStoryElementsNodeId converts echo context to params.
func (w *ServerInterfaceWrapper) GetStoryElementsNodeId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "nodeId" -------------
	var nodeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "nodeId", runtime.ParamLocationPath, ctx.Param("nodeId"), &nodeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter nodeId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetStoryElementsNodeId(ctx, nodeId)
	return err
}

// PutStoryElementsNodeId converts echo context to params.
func (w *ServerInterfaceWrapper) PutStoryElementsNodeId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "nodeId" -------------
	var nodeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "nodeId", runtime.ParamLocationPath, ctx.Param("nodeId"), &nodeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter nodeId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutStoryElementsNodeId(ctx, nodeId)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/player", wrapper.PostPlayer)
	router.GET(baseURL+"/player/:playerId", wrapper.GetPlayerPlayerId)
	router.PATCH(baseURL+"/player/:playerId", wrapper.PatchPlayerPlayerId)
	router.POST(baseURL+"/storyElements", wrapper.PostStoryElements)
	router.GET(baseURL+"/storyElements/:nodeId", wrapper.GetStoryElementsNodeId)
	router.PUT(baseURL+"/storyElements/:nodeId", wrapper.PutStoryElementsNodeId)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xVQW/bPAz9K4K+D9gliNPt5lvWDkOAYQta9DAMxaDaTKPCljSSSmYE+e+DZKeJa7tN",
	"BzTdKQ79JJLvkc8bmdnSWQOGSaYbSdkSShUf54WqAMOTQ+sAWUOM/9R5+OHKgUwlMWpzJ7cjCaXSRe8b",
	"YovVFSuuL9AMJXXvVch6oTJuYzqXNQGFqKrwP/OIYPh8qRzX5XaONJC5Qh6ub3bR+26tKbfli2raB+zt",
	"PWTcV7UnQKNKGMj5u65mYbFULFPpvc7l6DGyL9FV6OVTASUYPl67J7kf7iLbk97R8vryS2+qwa5XOgfb",
	"f6qv0cwabnrs5rA5DOhpHWtrXtinaybn5E3+5WQOddJNEULaLGw4nANlqCM/MpVTI6bzmVhYFEqcL60l",
	"EN+tR/FtbcQ0X4FhjyDuVAnjMJqai3DxMHI6n8mRXAFSneBsPBlPaknAKKdlKj/EUOCbl7GfxO1NyFLU",
	"ICigQo2zXKZybokboxpJhF8eiD/avArIgwlRzhU6i8eSewrpd1YXnv5HWMhU/pfsvTBpjDBpLo9Ehfs1",
	"Qi5TRg8xQM4aqsfh/eTsVbK2ZanfCApuKjIExZAL8lkGRAtfFNU4yky+LBVWQZCIEUoYWAt3cLoGNgwn",
	"m/p3lm9DbXfQw/VnaKieN9CoFKoSGJBk+mMjdSgxqCd3OyDdHtzmb3TAxXNGd9PhenJqrhEYNawgf0zw",
	"ZfNCqIbed/RAcBzlbNkztyH8pmy+7aqcXD7v8rAqj8W7juGOdOK2ErOLZkPo4JNKT1vRVQv6OjS3vvAn",
	"9qVu7jbl8b2AGvBye6LD433sJ5v4aX/apFoifI34o3bL7KDDm3VKXzqO61D0Ud5ED/Dal3zfCPuTk/cv",
	"7MfkjfbjOU/aSxb8SDM1nhTQgKudGh4Lmcols6M0SZTT48p6XMMtaYZxZku5vdn+CQAA//9DyxqM6Q0A",
	"AA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
