package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Patchable helpers for testability
var (
	decodeLastEvaluatedKeyTest func(string) (map[string]types.AttributeValue, error)
	encodeLastEvaluatedKeyTest func(map[string]types.AttributeValue) (string, error)
)

func TestListCategoriesForUser_Success(t *testing.T) {
	client := new(MockDynamoDBClient)
	tableName := "test-table"
	repo := &DynamoDbCategoriesRepository{client: client, tableName: tableName}

	input := models.ListCategoriesInput{
		UserID: "user1",
		Limit:  2,
	}

	items := []map[string]types.AttributeValue{
		{"category_name": &types.AttributeValueMemberS{Value: "cat1"}},
		{"category_name": &types.AttributeValueMemberS{Value: "cat2"}},
	}
	client.On("Query", mock.Anything, mock.Anything).
		Return(&dynamodb.QueryOutput{
			Items:            items,
			LastEvaluatedKey: map[string]types.AttributeValue{"category_name": &types.AttributeValueMemberS{Value: "cat2"}},
		}, nil)

	cats, nextKey, err := repo.ListCategoriesForUser(context.Background(), input)
	assert.NoError(t, err)
	assert.Len(t, cats, 2)
	assert.Equal(t, "cat1", cats[0].Name)
	assert.Equal(t, "cat2", cats[1].Name)
	assert.Equal(t, "cat2", nextKey)
}

func TestListCategoriesForUser_DynamoError(t *testing.T) {
	client := new(MockDynamoDBClient)
	repo := NewDynamoDbCategoriesRepository(client, "table")
	input := models.ListCategoriesInput{UserID: "user1"}
	client.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("dynamo error"))
	cats, nextKey, err := repo.ListCategoriesForUser(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, cats)
	assert.Equal(t, "", nextKey)
}

func TestListCategoriesForUser_UnmarshalError(t *testing.T) {
	client := new(MockDynamoDBClient)
	repo := NewDynamoDbCategoriesRepository(client, "table")
	input := models.ListCategoriesInput{UserID: "user1"}
	items := []map[string]types.AttributeValue{
		{"category_name": &types.AttributeValueMemberN{Value: "123"}},
	}
	client.On("Query", mock.Anything, mock.Anything).
		Return(&dynamodb.QueryOutput{Items: items}, nil)
	cats, nextKey, err := repo.ListCategoriesForUser(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, cats)
	assert.Equal(t, "", nextKey)
}

func TestListCategoriesForUser_StartKeyDecodeError(t *testing.T) {
	client := new(MockDynamoDBClient)
	repo := NewDynamoDbCategoriesRepository(client, "table")
	input := models.ListCategoriesInput{UserID: "user1", StartKey: "badkey"}
	decodeLastEvaluatedKeyTest = func(string) (map[string]types.AttributeValue, error) {
		return nil, errors.New("decode error")
	}
	cats, nextKey, err := repo.ListCategoriesForUser(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, cats)
	assert.Equal(t, "", nextKey)
	decodeLastEvaluatedKeyTest = nil
}

func TestListCategoriesForUser_NextKeyEncodeError(t *testing.T) {
	client := new(MockDynamoDBClient)
	repo := NewDynamoDbCategoriesRepository(client, "table")
	input := models.ListCategoriesInput{UserID: "user1"}
	items := []map[string]types.AttributeValue{
		{"category_name": &types.AttributeValueMemberS{Value: "cat1"}},
	}
	client.On("Query", mock.Anything, mock.Anything).
		Return(&dynamodb.QueryOutput{Items: items, LastEvaluatedKey: map[string]types.AttributeValue{"category_name": &types.AttributeValueMemberS{Value: "cat1"}}}, nil)
	encodeLastEvaluatedKeyTest = func(map[string]types.AttributeValue) (string, error) {
		return "", errors.New("encode error")
	}
	cats, nextKey, err := repo.ListCategoriesForUser(context.Background(), input)
	assert.Error(t, err)
	assert.Nil(t, cats)
	assert.Equal(t, "", nextKey)
	encodeLastEvaluatedKeyTest = nil
}
