package handler

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gaku3601/serverless-go/func/create/dynamo"
	"github.com/tidwall/gjson"
)

type (
	handler struct {
		dynamoModel dynamo.DynamoModel
	}
)

func NewHandler(d dynamo.DynamoModel) handler {
	return handler{d}
}

func (h handler) CreateData(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// jsonの値を取得
	title := gjson.Get(request.Body, "title").String()
	h.dynamoModel.Create(title)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("success\n"),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}
