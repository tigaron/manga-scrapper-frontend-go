package main

import (
	"manga-scraper-fe-go/pkg/handlers"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	seriesTable   = os.Getenv("SERIES_TABLE")
	chaptersTable = os.Getenv("CHAPTERS_TABLE")
	ddbClient     dynamodbiface.DynamoDBAPI
)

func main() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ddbClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "GET" {
		return handlers.UnhandledMethod()
	}

	switch request.Resource {
	case "/series/{seriesId}/chapters/{chaptersId}":
		return handlers.GetChaptersById(request, chaptersTable, ddbClient)
	case "/series/{seriesId}/chapters":
		return handlers.GetChaptersBySeries(request, chaptersTable, ddbClient)
	case "/series/{seriesId}":
		return handlers.GetSeriesById(request, seriesTable, ddbClient)
	case "/series":
		if _, exist := request.QueryStringParameters["provider"]; exist {
			return handlers.GetSeriesByProvider(request, seriesTable, ddbClient)
		} else {
			return handlers.GetAllSeries(request, seriesTable, ddbClient)
		}
	default:
		return handlers.UnhandledResource()
	}
}
