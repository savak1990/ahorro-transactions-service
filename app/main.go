package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/config"
	"github.com/savak1990/transactions-service/app/handler"
	"github.com/savak1990/transactions-service/app/repo"
	"github.com/savak1990/transactions-service/app/schema"
	"github.com/savak1990/transactions-service/app/service"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

// Lambda handler for API Gateway
func lambdaHandler(adapter *gorillamux.GorillaMuxAdapter) func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		var claims map[string]interface{}
		if event.RequestContext.Authorizer != nil && event.RequestContext.Authorizer.JWT != nil {
			claims = make(map[string]interface{})
			for k, v := range event.RequestContext.Authorizer.JWT.Claims {
				claims[k] = v
			}
		} else {
			claims = map[string]interface{}{}
		}
		log.WithFields(log.Fields{
			"method":      event.RequestContext.HTTP.Method,
			"path":        event.RawPath,
			"query":       event.RawQueryString,
			"headers":     event.Headers,
			"body":        event.Body,
			"cognito_sub": claims["sub"],
		}).Info("REQUEST")

		resp, err := adapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV2(&event))

		v2resp := resp.Version2()
		log.WithFields(log.Fields{
			"status":  v2resp.StatusCode,
			"headers": v2resp.Headers,
			"body":    v2resp.Body,
		}).Info("RESPONSE")

		return *v2resp, err
	}
}

func main() {
	// Validate embedded schema first
	if err := schema.ValidateSchemaEmbedded(); err != nil {
		log.WithError(err).Fatal("Failed to validate embedded OpenAPI schema")
	}

	appCfg := config.LoadConfig()
	log.WithFields(log.Fields{
		"region":  appCfg.AWSRegion,
		"profile": appCfg.AWSProfile,
		"db_host": appCfg.DBHost,
		"db_port": appCfg.DBPort,
		"db_name": appCfg.DBName,
	}).Info("Loaded config")

	// Initialize GORM database connection
	gormDB := aws.GetGormDB(appCfg)

	// Initialize PostgreSQL repositories
	repo := repo.NewPostgreSQLRepository(gormDB)
	service := service.NewServiceImpl(repo)
	serviceHandler := handler.NewHandlerImpl(service)

	commonHandler := handler.NewCommonHandlerImpl(gormDB)

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(handler.EnsureAwsRegionHeader(appCfg.AWSRegion))

	// Common APIs
	router.HandleFunc("/health", commonHandler.HandleHealth).Methods("GET")
	router.HandleFunc("/info", commonHandler.HandleInfo).Methods("GET")

	// Schema APIs (serve embedded OpenAPI schema)
	router.HandleFunc("/docs", schema.ServeSwaggerUIHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/schema", schema.ServeSwaggerUIHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/schema/raw", schema.ServeSchemaRawHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/schema/json", schema.ServeSchemaJSONHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/schema/info", schema.ServeSchemaInfoHandler()).Methods("GET", "OPTIONS")

	// Transactions APIs
	router.HandleFunc("/transactions", serviceHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions", serviceHandler.ListTransactions).Methods("GET")
	router.HandleFunc("/transactions/{transaction_id}", serviceHandler.GetTransaction).Methods("GET")
	router.HandleFunc("/transactions/{transaction_id}", serviceHandler.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{transaction_id}", serviceHandler.DeleteTransaction).Methods("DELETE")

	// Balances APIs
	router.HandleFunc("/balances", serviceHandler.CreateBalance).Methods("POST")
	router.HandleFunc("/balances", serviceHandler.ListBalances).Methods("GET")
	router.HandleFunc("/balances/{balance_id}", serviceHandler.GetBalance).Methods("GET")
	router.HandleFunc("/balances/{balance_id}", serviceHandler.UpdateBalance).Methods("PUT")
	router.HandleFunc("/balances/{balance_id}", serviceHandler.DeleteBalance).Methods("DELETE")

	// Categories APIs
	router.HandleFunc("/categories", serviceHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories", serviceHandler.ListCategories).Methods("GET")
	router.HandleFunc("/categories/{category_id}", serviceHandler.DeleteCategory).Methods("DELETE")

	// Lambda/API Gateway integration: use the muxadapter if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		adapter := gorillamux.New(router)
		lambda.Start(lambdaHandler(adapter))
		return
	}

	// Local dev server
	log.Info("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
