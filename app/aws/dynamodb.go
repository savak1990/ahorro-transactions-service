package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dbClient     *dynamodb.Client
	dbClientOnce sync.Once
)

func GetDynamoDbClient(cfg aws.Config) *dynamodb.Client {
	dbClientOnce.Do(func() {
		dbClient = dynamodb.NewFromConfig(cfg)
	})
	if dbClient == nil {
		panic("DynamoDB client is not initialized")
	}
	return dbClient
}
