package dynamo

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	DynamoModel interface {
		Create(title string)
		UpdateSequence(svc *dynamodb.DynamoDB, tableName string) *string
	}

	dynamoModel struct {
		svc *dynamodb.DynamoDB
	}
)

func NewDynamoModel(svc *dynamodb.DynamoDB) *dynamoModel {
	return &dynamoModel{svc: svc}
}

func (d *dynamoModel) Create(title string) {
	// 作成時間を取得
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t := time.Now().In(jst)
	lt := t.Format("20060102150405.000")

	id := d.UpdateSequence(d.svc, os.Getenv("SEQUENCE_TABLE"))

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

	_, putErr := d.svc.PutItem(putParams)
	if putErr != nil {
		panic(putErr)
	}
}

func (d *dynamoModel) UpdateSequence(svc *dynamodb.DynamoDB, tableName string) *string {
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
