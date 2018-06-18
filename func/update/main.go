package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/tidwall/gjson"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	// jsonの値を取得
	title := gjson.Get(request.Body, "title").String()
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	svc := dynamodb.New(sess)
	putParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(id),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"Title": {
				Value: &dynamodb.AttributeValue{
					S: aws.String(title),
				},
			},
		},
		// 返却内容を記載するのを忘れない！！！！
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	putItem, putErr := svc.UpdateItem(putParams)
	if putErr != nil {
		panic(fmt.Sprintf("error:%#v", putErr))
	}
	type Obj struct {
		Title string
	}
	obj := Obj{}
	dynamodbattribute.UnmarshalMap(putItem.Attributes, &obj)
	j, _ := json.Marshal(obj)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", string(j)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
