package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/savak1990/transactions-service/app/models"
)

type DynamoDbCategoriesRepository struct {
	client    DynamoDbClient
	tableName string
}

func NewDynamoDbCategoriesRepository(client DynamoDbClient, tableName string) *DynamoDbCategoriesRepository {
	return &DynamoDbCategoriesRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDbCategoriesRepository) ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	// Decode the start key if provided
	var startKey map[string]types.AttributeValue
	var err error
	if input.StartKey != "" {
		startKey, err = decodeLastEvaluatedKey(input.StartKey)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode start key: %w", err)
		}
	}

	// Build the DynamoDB query input
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: input.UserID},
		},
		ScanIndexForward:  aws.Bool(true), // Ascending order by sort key
		ExclusiveStartKey: startKey,
	}
	if input.Limit > 0 {
		queryInput.Limit = aws.Int32(int32(input.Limit))
	}

	resp, err := r.client.Query(ctx, queryInput)
	if err != nil {
		return nil, "", fmt.Errorf("dynamodb query failed: %w", err)
	}

	categories := make([]models.Category, 0, len(resp.Items))
	for _, item := range resp.Items {
		var cat models.Category
		if err := attributevalue.UnmarshalMap(item, &cat); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal category: %w", err)
		}
		categories = append(categories, cat)
	}

	nextKey, err := encodeLastEvaluatedKey(resp.LastEvaluatedKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode last evaluated key: %w", err)
	}

	return categories, nextKey, nil
}
