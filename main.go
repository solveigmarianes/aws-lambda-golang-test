package main

import (
	"context"
	"errors"
	"fmt"

	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	HTTPMethodNotSupported = errors.New("no name was provided in the HTTP body")
)

type CalculationRequest struct {
	A      int    `json:"a"`
	B      int    `json:"b"`
	Method string `json:"method"`
}

type Result struct {
	Result int `json:"result"`
}

func calculate(req CalculationRequest) Result {
	switch req.Method {
	case "ADD":
		return Result{req.A + req.B}
	case "SUBTRACT":
		return Result{req.A - req.B}
	case "MULTIPLY":
		return Result{req.A * req.B}
	case "DIVIDE":
		return Result{req.A / req.B}
	default:
		return Result{0}
	}
}

func createResponse(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var calculationRequest CalculationRequest
	err := json.Unmarshal([]byte(request.Body), &calculationRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Could not unmarshal Json",
		}, nil
	}

	response, respErr := json.Marshal(calculate(calculationRequest))
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
	fmt.Printf("Body size = %d. \n", len(request.Body))
	fmt.Println(request.Body)
	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}
	if request.HTTPMethod == "GET" {
		fmt.Printf("GET METHOD\n")
		return events.APIGatewayProxyResponse{Body: "GET", StatusCode: 200}, nil
	} else if request.HTTPMethod == "POST" {
		fmt.Printf("POST METHOD\n")
		return createResponse(request)
	} else {
		fmt.Printf("NEITHER\n")
		return events.APIGatewayProxyResponse{}, HTTPMethodNotSupported
	}
}

func main() {
	lambda.Start(HandleRequest)
}
