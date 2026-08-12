package internal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/middleware"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	api "github.com/Knoblauchpilze/galactic-sovereign/api"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swagv2 "github.com/swaggo/swag/v2"
)

func SwaggerEndpoints(config server.Config) ([]*rest.Route, error) {
	var out []*rest.Route

	doc, err := loadAndRewriteSwaggerToUICompatibleDoc()
	if err != nil {
		return out, err
	}

	handler := createStaticOpenApiSpecHandler(doc)
	staticGet := rest.NewRawRoute(http.MethodGet, "/openapi3.json", handler)
	out = append(out, staticGet)

	// The swagger UI needs to be configured to provide the URL of the
	// API spec file. Due to the rewriting, the spec content is served
	// by a dedicated endpoint. The final path of this endpoint needs
	// to include the base path which gets added when the route is
	// registered in the server.
	// The UI is reachable under: BASE_PATH/swagger/index.html
	swaggerUi := rest.NewRawRoute(
		http.MethodGet,
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			ginSwagger.URL(rest.ConcatenateEndpoints(config.BasePath, staticGet.Path())),
		),
	)
	out = append(out, swaggerUi)

	return out, nil
}

// createStaticOpenApiSpecHandler the swagger UI expects the API spec to be served at
// a specific URL. The setup works (almost) out of the box with the ginSwagger handler
// but does not handle well the re-writing necessary to work around the spec version.
// AI suggested to use a static endpoint serving the rewritten spec and to point the
// swagger UI handler to it.
// This function creates the static handler by defining an endpoint which serves the
// modified file as raw data.
func createStaticOpenApiSpecHandler(apiSpecDoc []byte) middleware.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", apiSpecDoc)
	}
}

// loadAndRewriteSwaggerToUICompatibleDoc the API spec generation step produces a spec
// in Open API spec 3.1 format. However, the swagger UI only accepts specs up to 3.0.X.
// This function allows to:
//   - load the specs from the (generated) spec file
//   - modify/hack the version indicated in it to be 3.0.3 (compatible with the UI)
//   - return the document as a string
func loadAndRewriteSwaggerToUICompatibleDoc() ([]byte, error) {
	doc, err := swagv2.ReadDoc(api.SwaggerInfo.InstanceName())
	if err != nil {
		return nil, errors.New("failed to read OpenAPI document")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(doc), &payload); err != nil {
		return nil, err
	}

	payload["openapi"] = "3.0.3"

	out, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return out, nil
}
