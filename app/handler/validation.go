package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/savak1990/transactions-service/app/schema"
	log "github.com/sirupsen/logrus"
)

// ValidationMiddleware provides OpenAPI schema validation for incoming requests
type ValidationMiddleware struct {
	router routers.Router
}

// NewValidationMiddleware creates a new validation middleware using the embedded OpenAPI schema
func NewValidationMiddleware() (*ValidationMiddleware, error) {
	// Load the embedded OpenAPI schema
	loader := &openapi3.Loader{Context: context.Background()}
	doc, err := loader.LoadFromData([]byte(schema.OpenAPISchema))
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI schema: %w", err)
	}

	// Validate the schema
	err = doc.Validate(context.Background())
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI schema: %w", err)
	}

	// Create router for request validation
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI router: %w", err)
	}

	return &ValidationMiddleware{
		router: router,
	}, nil
}

// ValidateRequest validates an HTTP request against the OpenAPI schema
func (v *ValidationMiddleware) ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip validation for CORS preflight requests
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Skip validation for certain endpoints that don't need it
		if v.shouldSkipValidation(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Find the route in the OpenAPI spec
		route, pathParams, err := v.router.FindRoute(r)
		if err != nil {
			// If route not found in OpenAPI spec, let it pass through
			// This allows for endpoints not defined in the schema to work
			log.WithField("path", r.URL.Path).WithField("method", r.Method).Debug("Route not found in OpenAPI spec, skipping validation")
			next.ServeHTTP(w, r)
			return
		}

		// Create request validation input
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
					// Skip authentication validation - we handle auth separately
					return nil
				},
			},
		}

		// Validate the request
		err = openapi3filter.ValidateRequest(context.Background(), requestValidationInput)
		if err != nil {
			v.handleValidationError(w, err)
			return
		}

		// Request is valid, proceed
		next.ServeHTTP(w, r)
	})
}

// shouldSkipValidation determines if validation should be skipped for certain paths
func (v *ValidationMiddleware) shouldSkipValidation(path string) bool {
	skipPaths := []string{
		"/health",
		"/info",
		"/docs",
		"/schema",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// handleValidationError formats and returns validation errors
func (v *ValidationMiddleware) handleValidationError(w http.ResponseWriter, err error) {
	log.WithError(err).Warn("Request validation failed")

	// Extract meaningful error message
	errorMessage := "Request validation failed"

	if requestErr, ok := err.(*openapi3filter.RequestError); ok {
		if requestErr.Parameter != nil {
			errorMessage = fmt.Sprintf("Invalid parameter '%s': %s", requestErr.Parameter.Name, requestErr.Reason)
		} else if requestErr.RequestBody != nil {
			errorMessage = fmt.Sprintf("Invalid request body: %s", requestErr.Reason)
		} else {
			errorMessage = requestErr.Reason
		}
	} else if securityErr, ok := err.(*openapi3filter.SecurityRequirementsError); ok {
		errorMessage = fmt.Sprintf("Security requirements not met: %s", securityErr.Error())
	} else {
		errorMessage = err.Error()
	}

	// Return standardized error response
	WriteJSONError(w, http.StatusBadRequest, "SchemaValidationError", errorMessage)
}

// ValidateResponse validates HTTP responses against the OpenAPI schema (optional)
func (v *ValidationMiddleware) ValidateResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, we'll just pass through without response validation
		// Response validation can be added later if needed
		next.ServeHTTP(w, r)
	})
}
