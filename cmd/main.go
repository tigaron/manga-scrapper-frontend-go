package main

import (
	"log"
	"manga-scraper-fe-go/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	tableName = os.Getenv("SERIES_TABLE")
	ddbClient dynamodbiface.DynamoDBAPI
)

func main() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ddbClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("Request HTTP method: %v", request.HTTPMethod)
	if request.HTTPMethod != "GET" {
		return handlers.UnhandledMethod()
	}

	_, exist := request.PathParameters["seriesId"]
	if exist {
		return handlers.GetSeriesById(request, tableName, ddbClient)
	}

	_, exist = request.QueryStringParameters["provider"]
	if exist {
		return handlers.GetSeriesByProvider(request, tableName, ddbClient)
	}

	return handlers.GetAllSeries(request, tableName, ddbClient)
}
