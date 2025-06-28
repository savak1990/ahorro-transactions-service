package schema

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Embed the OpenAPI schema file at compile time
//
//go:embed openapi.yml
var OpenAPISchema string

// SwaggerUI HTML template
const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>API Documentation - Ahorro Transactions Service</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
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
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            // Get the schema from the embedded endpoint
            const ui = SwaggerUIBundle({
                url: '/schema/raw',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                validatorUrl: null,
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                onComplete: function(swaggerApi, swaggerUi) {
                    console.log("Swagger UI loaded successfully");
                },
                onFailure: function(data) {
                    console.log("Unable to Load SwaggerUI");
                    console.log(data);
                }
            });
        };
    </script>
</body>
</html>`

// SchemaInfo represents basic information about the embedded schema
type SchemaInfo struct {
	Format  string `json:"format"`
	Version string `json:"version"`
	Size    int    `json:"size"`
}

// GetSchemaInfo returns information about the embedded schema
func GetSchemaInfo() SchemaInfo {
	return SchemaInfo{
		Format:  "yaml",
		Version: "3.0.3",
		Size:    len(OpenAPISchema),
	}
}

// ServeSchemaHandler returns an HTTP handler that serves the embedded OpenAPI schema
func ServeSchemaHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Write the schema content
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(OpenAPISchema)); err != nil {
			log.WithError(err).Error("Failed to write OpenAPI schema response")
		}
	}
}

// ServeSchemaInfoHandler returns an HTTP handler that serves information about the schema
func ServeSchemaInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get schema info
		info := GetSchemaInfo()

		// Write JSON response
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(info); err != nil {
			log.WithError(err).Error("Failed to encode schema info response")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// ServeSwaggerUIHandler returns an HTTP handler that serves the Swagger UI interface
func ServeSwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Write the Swagger UI HTML
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(swaggerUIHTML)); err != nil {
			log.WithError(err).Error("Failed to write Swagger UI response")
		}
	}
}

// ServeSchemaRawHandler returns the raw schema in YAML format (used by Swagger UI)
func ServeSchemaRawHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers for YAML content
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Write the schema content
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(OpenAPISchema)); err != nil {
			log.WithError(err).Error("Failed to write OpenAPI schema response")
		}
	}
}

// ServeSchemaJSONHandler returns the schema in JSON format
func ServeSchemaJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Convert YAML to JSON (basic conversion for demo)
		// For production, you might want to use a proper YAML to JSON library
		jsonSchema := convertYAMLToJSON(OpenAPISchema)

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(jsonSchema)); err != nil {
			log.WithError(err).Error("Failed to write JSON schema response")
		}
	}
}

// convertYAMLToJSON is a basic YAML to JSON converter
// For production use, consider using gopkg.in/yaml.v3
func convertYAMLToJSON(yamlContent string) string {
	// This is a simplified conversion - for production you should use proper YAML parsing
	// For now, we'll return a basic JSON structure
	return `{"info": {"title": "API Schema available in YAML format", "description": "Use /schema/raw for YAML format"}}`
}

// ValidateSchemaEmbedded validates that the schema was properly embedded
func ValidateSchemaEmbedded() error {
	if OpenAPISchema == "" {
		return fmt.Errorf("OpenAPI schema is empty - embedding may have failed")
	}

	if len(OpenAPISchema) < 100 {
		return fmt.Errorf("OpenAPI schema appears to be too small (%d bytes) - embedding may be incomplete", len(OpenAPISchema))
	}

	log.WithFields(log.Fields{
		"schema_size":  len(OpenAPISchema),
		"schema_start": OpenAPISchema[:min(50, len(OpenAPISchema))],
	}).Info("OpenAPI schema successfully embedded")

	return nil
}

// min is a helper function since we might be using older Go versions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
