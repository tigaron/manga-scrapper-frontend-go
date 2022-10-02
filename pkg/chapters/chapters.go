package chapters

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

type Chapters struct {
	ChaptersProvider string   `json:"_type"`
	ChaptersId       string   `json:"_id"`
	ChaptersTitle    string   `json:"ChapterTitle"`
	ChaptersNumber   string   `json:"ChapterNumber"`
	ChaptersDate     string   `json:"ChapterDate"`
	ChaptersUrl      string   `json:"ChapterUrl"`
	ChaptersShortUrl string   `json:"ChapterShortUrl"`
	ChaptersOrder    int      `json:"ChapterOrder"`
	ChaptersPrev     string   `json:"ChapterPrevSlug"`
	ChaptersNext     string   `json:"ChapterNextSlug"`
	ChaptersContent  []string `json:"ChapterContent"`
	ScrapeDate       string   `json:"ScrapeDate"`
}

func FetchChaptersBySeries(provider string, seriesId string, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Chapters, error) {
	keyCond := expression.Key("_type").Equal(expression.Value(provider + "_" + seriesId))
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
		log.Printf("Couldn't get result of '%v_%v'. Here's why: %v\n", provider, seriesId, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Chapters)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchChaptersBySeriesPaginated(provider string, seriesId string, pageSize int64, pageNum int, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*[]Chapters, error) {
	keyCond := expression.Key("_type").Equal(expression.Value(provider + "_" + seriesId))
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
		log.Printf("Couldn't get result of '%v_%v'. Here's why: %v\n", provider, seriesId, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]Chapters)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchOneChapters(provider string, seriesId string, chaptersId string, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*Chapters, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_type": {
				S: aws.String(provider + "_" + seriesId),
			},
			"_id": {
				S: aws.String(chaptersId),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := ddbClient.GetItem(input)
	if err != nil {
		log.Printf("Couldn't get result of '%v' in '%v_%v'. Here's why: %v\n", chaptersId, seriesId, provider, err)
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(Chapters)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Printf("Couldn't unmarshal result. Here's why: %v\n", err)
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}
