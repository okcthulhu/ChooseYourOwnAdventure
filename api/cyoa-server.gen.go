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
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Retrieve a player's state.
	// (GET /player/{playerId})
	GetPlayerPlayerId(ctx echo.Context, playerId string) error
	// Retrieve a story node by its ID.
	// (GET /storyElements/{nodeId})
	GetStoryElementsNodeId(ctx echo.Context, nodeId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetPlayerPlayerId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPlayerPlayerId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "playerId" -------------
	var playerId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "playerId", runtime.ParamLocationPath, ctx.Param("playerId"), &playerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter playerId: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPlayerPlayerId(ctx, playerId)
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

	router.GET(baseURL+"/player/:playerId", wrapper.GetPlayerPlayerId)
	router.GET(baseURL+"/storyElements/:nodeId", wrapper.GetStoryElementsNodeId)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xUT2/TThD9Kqv5/SQukR3g5ltpEaqEIGrhgFAP2/Uk3mLvLjPjRFaU7452N5ZT4lbl",
	"grjZ82/fvDczezC+C96hE4ZqD2wa7HT6XLV6QLoVLRh/A/mAJBaTU5PYtTY5yQp26UOGgFABC1m3gcNi",
	"NGgiPcR/0xOhk8tGB0GaTTmGrDTJrD8kWNf1rHNnufbdH4GaDP7+AY3EiFvxNLxvsUMn572bCf0ZKV9v",
	"Ps6+6XSHs46trdHPZ80hM97JEdT5G77GJ3jxQax3s7xMtX8XKxwl+OtNPqPiU2jPy0STdWsfk2tkQzZx",
	"ABVcOHWxulZrT0qry8Z7RvXN96Q+75y6qLfopCdUG92hWnvTM9bKOyWoTWPdRl29+6L4h21bLmABYqWN",
	"7z5daGeliUmwgC0SZxCvi2WxzNKg08FCBW+TKfIuTeq5zLNe7seZP0TrBpMqURMdO4qSwweUvK+rcT2S",
	"frpDQWKovu/BxmdjbRiVmnZpAYQ/e0tYQyXU4+J4COYEu4vBHLzjPBBvlsu0FdNk6hBaaxK48oFju/uT",
	"ev8TrqGC/8rp8pTHs1Oe3pyk4GPlsltx9CtCIYtbrIskP/ddp2mACm6ODqVVbvAV55QcWPLJdnO5z2vz",
	"LLWn94A/5TV7Cb9uDP032H101mboTX4VQb+IXJ7C7wdlhdX1VZHrMtJ25KWnFipoRAJXZamDLQbf0w7v",
	"2QoWxndwuDv8CgAA//9lcWfiiAYAAA==",
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