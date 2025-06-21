package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/savak1990/transactions-service/app/models"
)

type DynamoDbTransactionsRepository struct {
	client    DynamoDbClient
	tableName string
}

func NewDynamoDbTransactionsRepository(client DynamoDbClient, tableName string) *DynamoDbTransactionsRepository {
	return &DynamoDbTransactionsRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDbTransactionsRepository) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	item, err := attributevalue.MarshalMap(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(transaction_id)"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put transaction: %w", err)
	}
	return &tx, nil
}

func (r *DynamoDbTransactionsRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: filter.UserID},
		},
		ScanIndexForward: aws.Bool(false),
	}
	if filter.Count > 0 {
		queryInput.Limit = aws.Int32(int32(filter.Count))
	}
	if filter.StartKey != "" {
		startKey, err := decodeLastEvaluatedKey(filter.StartKey)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode start key: %w", err)
		}
		queryInput.ExclusiveStartKey = startKey
	}
	resp, err := r.client.Query(ctx, queryInput)
	if err != nil {
		return nil, "", fmt.Errorf("dynamodb query failed: %w", err)
	}
	transactions := make([]models.Transaction, 0, len(resp.Items))
	for _, item := range resp.Items {
		var tx models.Transaction
		if err := attributevalue.UnmarshalMap(item, &tx); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	nextKey, err := encodeLastEvaluatedKey(resp.LastEvaluatedKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode last evaluated key: %w", err)
	}
	return transactions, nextKey, nil
}

func (r *DynamoDbTransactionsRepository) GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":        &types.AttributeValueMemberS{Value: userID},
			"transaction_id": &types.AttributeValueMemberS{Value: transactionID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if result.Item == nil {
		return nil, nil
	}
	var tx models.Transaction
	if err := attributevalue.UnmarshalMap(result.Item, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return &tx, nil
}

func (r *DynamoDbTransactionsRepository) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	item, err := attributevalue.MarshalMap(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}
	return &tx, nil
}

func (r *DynamoDbTransactionsRepository) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id":        &types.AttributeValueMemberS{Value: userID},
			"transaction_id": &types.AttributeValueMemberS{Value: transactionID},
		},
		ConditionExpression: aws.String("attribute_not_exists(transaction_id)"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

var _ TransactionsRepo = (*DynamoDbTransactionsRepository)(nil)
