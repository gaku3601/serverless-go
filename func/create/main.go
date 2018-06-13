package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 一意なidを生成
	guid := xid.New()
	title := gjson.Get(request.Body, "title").String()
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(guid.String()),
			},
			"Title": {
				S: aws.String(title),
			},
		},
	}

	_, putErr := svc.PutItem(putParams)
	if putErr != nil {
		panic(putErr)
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Insert DynamoDB: ID: %v, Title: %v \n", guid.String(), title),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
