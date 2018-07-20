package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gaku3601/serverless-go/func/create/dynamo"
	"github.com/gaku3601/serverless-go/func/create/handler"
)

func main() {
	svc := dynamodb.New(session.New(), &aws.Config{})
	d := dynamo.NewDynamoModel(svc)
	h := handler.NewHandler(d)
	lambda.Start(h.Router)
}
