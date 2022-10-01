package series

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToBuildExpression = "failed to build expression"
)

type Series struct {
	SeriesProvider string `json:"_type"`
	SeriesId       string `json:"_id"`
	SeriesTitle    string `json:"MangaTitle"`
	SeriesCover    string `json:"MangaCover"`
	SeriesUrl      string `json:"MangaUrl"`
	SeriesShortUrl string `json:"MangaShortUrl"`
	SeriesSynopsis string `json:"MangaSynopsis"`
	ScrapeDate     string `json:"ScrapeDate"`
}

func FetchAllSeries(tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := ddbClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchSeriesByProvider(provider string, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	keyCond := expression.KeyEqual(expression.Key("_type"), expression.Value(provider))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, errors.New(ErrorFailedToBuildExpression)
	}

	input := &dynamodb.QueryInput{
		KeyConditionExpression: expr.KeyCondition(),
	}

	result, err := ddbClient.Query(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchOneSeries(provider string, seriesId string, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*Series, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_type": {
				S: aws.String(provider),
			},
			"_id": {
				S: aws.String(seriesId),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := ddbClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Series)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}
