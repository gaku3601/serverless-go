package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/tidwall/gjson"
)

func updateSequence(svc *dynamodb.DynamoDB, tableName string) *string {
	putParams := &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("SEQUENCE_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"TableName": {
				S: aws.String(tableName),
			},
		},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			"CurrentNumber": {
				Value: &dynamodb.AttributeValue{
					N: aws.String("1"),
				},
				Action: aws.String("ADD"),
			},
		},
		// 返却内容を記載するのを忘れない！！！！
		ReturnValues: aws.String("UPDATED_NEW"),
	}
	putItem, putErr := svc.UpdateItem(putParams)
	if putErr != nil {
		panic(fmt.Sprintf("error:%#v", putErr))
	}
	return putItem.Attributes["CurrentNumber"].N
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 作成時間を取得
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t := time.Now().In(jst)
	lt := t.Format("20060102150405.000")

	// jsonの値を取得
	title := gjson.Get(request.Body, "title").String()
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	id := updateSequence(svc, os.Getenv("SEQUENCE_TABLE"))

	putParams := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: id,
			},
			"Title": {
				S: aws.String(title),
			},
			"CreateDate": {
				S: aws.String(lt),
			},
		},
	}

	_, putErr := svc.PutItem(putParams)
	if putErr != nil {
		panic(putErr)
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("[Insert DynamoDB] ID: %v, Title: %v, CreateDate: %v \n", *id, title, lt),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
