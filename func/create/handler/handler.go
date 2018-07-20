package handler

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gaku3601/serverless-go/func/create/dynamo"
	"github.com/gaku3601/serverless-go/func/create/gombda"
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

func (h handler) Router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	g := gombda.New(request)
	g.POST("/func", h.Create)
	g.GET("/func/{id}", h.Show)
	g.DELETE("/func/{id}", h.Destroy)
	g.PATCH("/func/{id}", h.Update)

	return g.Start()
}

func (h handler) Create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

func (h handler) Show(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	j := h.dynamoModel.Show(id)

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("%v", string(j)),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: 200,
	}, nil
}

func (h handler) Destroy(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	h.dynamoModel.Destroy(id)

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("Delete Success"),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: 200,
	}, nil
}

func (h handler) Update(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	// jsonの値を取得
	title := gjson.Get(request.Body, "title").String()
	j := h.dynamoModel.Update(id, title)

	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("%v", string(j)),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: 200,
	}, nil
}
