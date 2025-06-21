package repo

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
)

var dynamoDbMarshalMapFunc = attributevalue.MarshalMap
var dynamoDbMarshalMapJsonFunc = attributevalue.MarshalMapJSON
var dynamoDbUnmarshalMapFunc = attributevalue.UnmarshalMap
var dynamoDbUnmarshalMapJsonFunc = attributevalue.UnmarshalMapJSON
var dynamoDbUnmarshalListOfMapsFunc = attributevalue.UnmarshalListOfMaps

func encodeLastEvaluatedKey(key map[string]types.AttributeValue) (string, error) {
	if len(key) == 0 {
		return "", nil
	}

	b, err := dynamoDbMarshalMapJsonFunc(key)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func decodeLastEvaluatedKey(token string) (map[string]types.AttributeValue, error) {
	if token == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.WithError(err).Error("Failed to decode last evaluated key")
		return nil, err
	}

	key, err := dynamoDbUnmarshalMapJsonFunc(b)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal last evaluated key")
		return nil, err
	}
	return key, nil
}
