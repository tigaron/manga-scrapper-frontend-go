package handlers

import (
	"log"
	"manga-scraper-fe-go/pkg/series"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorMethodNotAllowed     = "method not allowed"
	ErrorInvalidLimitValue    = "invalid limit value"
	ErrorInvalidProviderValue = "invalid provider value"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetAllSeries(request events.APIGatewayProxyRequest, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	limit, exist := request.QueryStringParameters["limit"]

	if exist { // If 'limit' is provided --> fetch all Series data with pagination
		log.Printf("Derived limit from query: %v", limit)
		pageSize, err := strconv.ParseInt(limit, 10, 64)
		if pageSize == 0 || err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorInvalidLimitValue)})
		}

		var pageNum int
		page, exist := request.QueryStringParameters["page"]
		log.Printf("Derived page from query: %v", page)
		if exist { // If 'page' is provided --> convert value to integer
			pageNum, _ = strconv.Atoi(page)
		} else { // Otherwise --> set value to 1
			pageNum = 1
		}

		result, err := series.FetchAllSeriesPaginated(pageSize, pageNum, tableName, ddbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return apiResponse(http.StatusOK, result)
	} else { // Otherwise --> fetch all Series data without pagination
		result, err := series.FetchAllSeries(tableName, ddbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return apiResponse(http.StatusOK, result)
	}
}

func GetSeriesByProvider(request events.APIGatewayProxyRequest, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	provider, exist := request.QueryStringParameters["provider"]
	if !exist {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorInvalidProviderValue)})
	}

	log.Printf("Derived provider from query: %v", provider)

	limit, exist := request.QueryStringParameters["limit"]
	if exist { // If 'limit' is provided --> fetch Series data by provider with pagination
		log.Printf("Derived limit from query: %v", limit)
		pageSize, err := strconv.ParseInt(limit, 10, 64)
		if pageSize == 0 || err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorInvalidLimitValue)})
		}

		var pageNum int
		page, exist := request.QueryStringParameters["page"]
		log.Printf("Derived page from query: %v", page)
		if exist { // If 'page' is provided --> convert value to integer
			pageNum, _ = strconv.Atoi(page)
		} else { // Otherwise --> set value to 1
			pageNum = 1
		}

		result, err := series.FetchSeriesByProviderPaginated(provider, pageSize, pageNum, tableName, ddbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return apiResponse(http.StatusOK, result)
	} else { // Otherwise --> fetch Series data by provider without pagination
		result, err := series.FetchSeriesByProvider(provider, tableName, ddbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}

		return apiResponse(http.StatusOK, result)
	}
}

func GetSeriesById(request events.APIGatewayProxyRequest, tableName string, ddbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	seriesId := request.PathParameters["seriesId"]
	log.Printf("Derived seriesId from path: %v", seriesId)
	provider, exist := request.QueryStringParameters["provider"]
	if !exist {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorInvalidProviderValue)})
	}

	log.Printf("Derived provider from query: %v", provider)

	result, err := series.FetchOneSeries(provider, seriesId, tableName, ddbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return apiResponse(http.StatusOK, result)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
