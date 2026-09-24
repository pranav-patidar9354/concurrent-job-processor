package docs

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.json
var openAPISpec []byte

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Concurrent Job Processor - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
  <link rel="icon" type="image/png" href="https://fastapi.tiangolo.com/img/favicon.png">
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin: 0;
      background: #fafafa;
    }
    .topbar {
      display: none;
    }
    .scheme-container {
      display: none !important;
    }
    .swagger-ui .info {
      margin: 25px 0;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>
`

// RegisterDocsRoutes registers the OpenAPI spec and Swagger UI endpoints.
func RegisterDocsRoutes(router *gin.Engine) {
	router.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", openAPISpec)
	})

	docsHandler := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	}

	router.GET("/docs", docsHandler)
	router.GET("/swagger", docsHandler)
}
