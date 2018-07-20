package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	testSvc = dynamodb.New(session.New(), &aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region:   aws.String("ap-northeast-1"),
	})
)

func TestCreate(t *testing.T) {
	initProcess()
	defer tearDownProcess()

	d := NewDynamoModel(testSvc)
	d.Create("aaaa")
}

func initProcess() {
	//tableの作成
	createDynamoDataTable(testSvc)
	createSequenceTable(testSvc)
}

func tearDownProcess() {
	//tableの削除
	deleteSequenceTable(testSvc)
	deleteDynamoDataTable(testSvc)
}

func createDynamoDataTable(svc *dynamodb.DynamoDB) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("DynamoDataTable"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}

	svc.CreateTable(params)
}

func createSequenceTable(svc *dynamodb.DynamoDB) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String("SequenceTable"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("TableName"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("TableName"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}

	svc.CreateTable(params)
}

func deleteDynamoDataTable(svc *dynamodb.DynamoDB) {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String("DynamoDataTable"),
	}

	svc.DeleteTable(params)
}

func deleteSequenceTable(svc *dynamodb.DynamoDB) {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String("SequenceTable"),
	}

	svc.DeleteTable(params)
}
