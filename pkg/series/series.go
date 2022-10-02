package series

import (
	"errors"
	"log"

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
	SeriesProvider *string `json:"_type"`
	SeriesId       *string `json:"_id"`
	SeriesTitle    *string `json:"MangaTitle"`
	SeriesCover    *string `json:"MangaCover"`
	SeriesUrl      *string `json:"MangaUrl"`
	SeriesShortUrl *string `json:"MangaShortUrl"`
	SeriesSynopsis *string `json:"MangaSynopsis"`
	ScrapeDate     *string `json:"ScrapeDate"`
}

func FetchAllSeries(tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := ddbClient.Scan(input)
	if err != nil {
		log.Printf("Couldn't get any result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchAllSeriesPaginated(pageSize int64, pageNum int, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int64(pageSize),
	}

	result := new(dynamodb.ScanOutput)
	index := 0
	err := ddbClient.ScanPages(input, func(page *dynamodb.ScanOutput, lastPage bool) bool {
		index++
		if index == pageNum {
			result = page
		}

		return index != pageNum
	})

	if err != nil {
		log.Printf("Couldn't get any result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchSeriesByProvider(provider string, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	keyCond := expression.Key("_type").Equal(expression.Value(provider))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, errors.New(ErrorFailedToBuildExpression)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := ddbClient.Query(input)
	if err != nil {
		log.Printf("Couldn't get result of '%v'. Here's why: %v\n", provider, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchSeriesByProviderPaginated(provider string, pageSize int64, pageNum int, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Series, error) {
	keyCond := expression.Key("_type").Equal(expression.Value(provider))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, errors.New(ErrorFailedToBuildExpression)
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int64(pageSize),
	}

	result := new(dynamodb.QueryOutput)
	index := 0
	err = ddbClient.QueryPages(input, func(page *dynamodb.QueryOutput, lastPage bool) bool {
		index++
		if index == pageNum {
			result = page
		}

		return index != pageNum
	})

	if err != nil {
		log.Printf("Couldn't get result of '%v'. Here's why: %v\n", provider, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Series)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
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
		log.Printf("Couldn't get result of '%v' in '%v'. Here's why: %v\n", seriesId, provider, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Series)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}
