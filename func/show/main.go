package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v %v %v", request.QueryStringParameters["foo"], request.QueryStringParameters["bar"], os.Getenv("DYNAMODATATABLE")),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
