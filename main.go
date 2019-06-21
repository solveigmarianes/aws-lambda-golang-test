package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	HTTPMethodNotSupported               = errors.New("no name was provided in the HTTP body")
	calculationMethodNotSupportedMessage = "Invalid method. Valid methods are \"ADD\", \"SUBTRACT\", \"MULTIPLY\" and \"DIVIDE\""
)

type CalculationRequest struct {
	A      int    `json:"a"`
	B      int    `json:"b"`
	Method string `json:"method"`
}

type Result struct {
	Result int `json:"result"`
}

func calculate(req CalculationRequest) (Result, error) {
	switch req.Method {
	case "ADD":
		return Result{req.A + req.B}, nil
	case "SUBTRACT":
		return Result{req.A - req.B}, nil
	case "MULTIPLY":
		return Result{req.A * req.B}, nil
	case "DIVIDE":
		return Result{req.A / req.B}, nil
	default:
		return Result{0}, errors.New(calculationMethodNotSupportedMessage)
	}
}

func createResponse(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calculationRequest CalculationRequest
	err := json.Unmarshal([]byte(request.Body), &calculationRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Could not unmarshal Json, please provide a valid request",
		}, nil
	}

	result, resultErr := calculate(calculationRequest)
	if resultErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       calculationMethodNotSupportedMessage,
		}, nil
	}

	response, respErr := json.Marshal(result)
	if respErr != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Could not marshal Json",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" {
		return createResponse(request)
	} else {
		return events.APIGatewayProxyResponse{}, HTTPMethodNotSupported
	}
}

func main() {
	lambda.Start(HandleRequest)
}
