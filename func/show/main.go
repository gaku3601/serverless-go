package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// クエリパラメータの場合
	//id := request.QueryStringParameters["id"]
	// パスパラメータの場合
	id := request.PathParameters["id"]
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	params := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
		AttributesToGet: []*string{
			aws.String("ID"),
			aws.String("Title"),
		},
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	}

	resp, err := svc.GetItem(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("ID: %v Title: %v", *resp.Item["ID"].S, *resp.Item["Title"].S),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
