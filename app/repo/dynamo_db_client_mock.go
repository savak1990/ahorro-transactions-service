package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

type MockDynamoDBClient struct {
	mock.Mock
}

func NewMockDynamoDBClient() *MockDynamoDBClient {
	return &MockDynamoDBClient{}
}

func (m *MockDynamoDBClient) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {

	args := m.Called(ctx, params)
	var out *dynamodb.DescribeTableOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.DescribeTableOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context,
	input *dynamodb.PutItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {

	args := m.Called(ctx, input)
	var out *dynamodb.PutItemOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.PutItemOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context,
	input *dynamodb.GetItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {

	args := m.Called(ctx, input)
	var out *dynamodb.GetItemOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.GetItemOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) Query(ctx context.Context,
	params *dynamodb.QueryInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {

	args := m.Called(ctx, params)
	var out *dynamodb.QueryOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.QueryOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) UpdateItem(ctx context.Context,
	input *dynamodb.UpdateItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {

	args := m.Called(ctx, input)
	var out *dynamodb.UpdateItemOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.UpdateItemOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItem(ctx context.Context,
	input *dynamodb.DeleteItemInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {

	args := m.Called(ctx, input)
	var out *dynamodb.DeleteItemOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.DeleteItemOutput)
	}
	return out, args.Error(1)
}

func (m *MockDynamoDBClient) Scan(ctx context.Context,
	input *dynamodb.ScanInput,
	optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {

	args := m.Called(ctx, input)
	var out *dynamodb.ScanOutput
	if v := args.Get(0); v != nil {
		out = v.(*dynamodb.ScanOutput)
	}
	return out, args.Error(1)
}

// Ensure MockDynamoDBClient implements DynamoDbClient interface
var _ DynamoDbClient = (*MockDynamoDBClient)(nil)
