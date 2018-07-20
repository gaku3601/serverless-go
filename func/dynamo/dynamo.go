package dynamo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type (
	DynamoModel interface {
		Create(title string)
		Show(id string) string
		Destroy(id string)
		Update(id string, title string) string
		Index(startID string, endID string) string
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

func (d *dynamoModel) Show(id string) string {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(id),
			},
		},
		AttributesToGet: []*string{
			aws.String("ID"),
			aws.String("Title"),
		},
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	}

	resp, err := d.svc.GetItem(params)

	if err != nil {
		fmt.Println(err.Error())
	}
	type Obj struct {
		ID    string
		Title string
	}

	obj := Obj{}
	dynamodbattribute.UnmarshalMap(resp.Item, &obj)
	j, _ := json.Marshal(obj)
	return string(j)
}

func (d *dynamoModel) Destroy(id string) {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				N: aws.String(id),
			},
		},

		ReturnConsumedCapacity:      aws.String("NONE"),
		ReturnItemCollectionMetrics: aws.String("NONE"),
		ReturnValues:                aws.String("NONE"),
	}

	_, err := d.svc.DeleteItem(params)
	if err != nil {
		panic(fmt.Sprintf("error:%#v", err))
	}
}

func (d *dynamoModel) Update(id string, title string) string {
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

	putItem, putErr := d.svc.UpdateItem(putParams)
	if putErr != nil {
		panic(fmt.Sprintf("error:%#v", putErr))
	}
	type Obj struct {
		Title string
	}
	obj := Obj{}
	dynamodbattribute.UnmarshalMap(putItem.Attributes, &obj)
	j, _ := json.Marshal(obj)
	return string(j)
}

func (d *dynamoModel) Index(startID string, endID string) string {
	params := &dynamodb.ScanInput{
		TableName:            aws.String(os.Getenv("DYNAMO_DATA_TABLE")),
		ProjectionExpression: aws.String("ID, Title, CreateDate"),
		FilterExpression:     aws.String("#key BETWEEN :startKey AND :endKey"),
		ExpressionAttributeNames: map[string]*string{
			"#key": aws.String("ID"), // 項目名をプレースホルダに入れる
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":startKey": {
				N: &startID,
			},
			":endKey": {
				N: &endID,
			},
		},
	}
	resp, err := d.svc.Scan(params)
	if err != nil {
		fmt.Errorf("failed to make Query API call, %v", err)
	}

	type Data struct {
		ID         int
		Title      string
		CreateDate string
	}

	obj := []Data{}
	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &obj)
	sort.Slice(obj, func(i, j int) bool {
		return obj[i].ID > obj[j].ID
	})
	j, _ := json.Marshal(obj)
	return string(j)
}
