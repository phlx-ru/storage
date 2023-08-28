// Package storage provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	externalRef0 "storage/schema/common"
	externalRef1 "storage/schema/storage"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

const (
	IntegrationsScopes = "integrations.Scopes"
	JwtScopes          = "jwt.Scopes"
)

// UploadParams defines parameters for Upload.
type UploadParams struct {
	// Filename Filename
	Filename externalRef0.Filename `form:"filename" json:"filename" validate:"required,min=3,max=255"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/1/download/{uid})
	Download(c *gin.Context, uid externalRef1.Uid)

	// (OPTIONS /api/1/download/{uid})
	DownloadOptions(c *gin.Context, uid externalRef1.Uid)

	// (GET /api/1/files/list)
	FilesList(c *gin.Context)

	// (POST /api/1/upload)
	Upload(c *gin.Context, params UploadParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// Download operation middleware
func (siw *ServerInterfaceWrapper) Download(c *gin.Context) {

	var err error

	// ------------- Path parameter "uid" -------------
	var uid externalRef1.Uid

	err = runtime.BindStyledParameter("simple", false, "uid", c.Param("uid"), &uid)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter uid: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Download(c, uid)
}

// DownloadOptions operation middleware
func (siw *ServerInterfaceWrapper) DownloadOptions(c *gin.Context) {

	var err error

	// ------------- Path parameter "uid" -------------
	var uid externalRef1.Uid

	err = runtime.BindStyledParameter("simple", false, "uid", c.Param("uid"), &uid)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter uid: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DownloadOptions(c, uid)
}

// FilesList operation middleware
func (siw *ServerInterfaceWrapper) FilesList(c *gin.Context) {

	c.Set(IntegrationsScopes, []string{})

	c.Set(JwtScopes, []string{})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.FilesList(c)
}

// Upload operation middleware
func (siw *ServerInterfaceWrapper) Upload(c *gin.Context) {

	var err error

	c.Set(IntegrationsScopes, []string{})

	c.Set(JwtScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params UploadParams

	// ------------- Required query parameter "filename" -------------

	if paramValue := c.Query("filename"); paramValue != "" {

	} else {
		siw.ErrorHandler(c, fmt.Errorf("Query argument filename is required, but not found"), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "filename", c.Request.URL.Query(), &params.Filename)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter filename: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Upload(c, params)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/api/1/download/:uid", wrapper.Download)
	router.OPTIONS(options.BaseURL+"/api/1/download/:uid", wrapper.DownloadOptions)
	router.GET(options.BaseURL+"/api/1/files/list", wrapper.FilesList)
	router.POST(options.BaseURL+"/api/1/upload", wrapper.Upload)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"Lln82xvt/wQAAP///W6MkEAeAAA=",
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

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./common/schema.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	for rawPath, rawFunc := range externalRef1.PathToRawSpec(path.Join(pathPrefix, "./storage/schema.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
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
