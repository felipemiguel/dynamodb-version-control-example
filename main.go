package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	ID      string `json:"id"`
	Data    string `json:"data"`
	Version int    `json:"version"`
}

func main() {
	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		return endpoints.ResolvedEndpoint{
			URL: "http://localhost:4566",
		}, nil

	}
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		os.Exit(1)
	}

	svc := dynamodb.New(sess)
	tableName := "TableExample"
	id := "RecordID"
	initialVersion := 1
	initialData := "InitialData"

	err = insertRecord(svc, tableName, id, initialData, initialVersion)
	if err != nil {
		fmt.Println("Error inserting record:", err)
		return
	}

	newData := "NewData"
	expectedVersion := initialVersion

	err = updateRecordWithVersionCheck(svc, tableName, id, newData, expectedVersion)
	if err != nil {
		fmt.Println("Error updating record:", err)
		return
	}

	fmt.Println("Record updated successfully!")
}

func insertRecord(svc *dynamodb.DynamoDB, tableName, id, data string, version int) error {
	item := Item{
		ID:      id,
		Data:    data,
		Version: version,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	return err
}

func updateRecordWithVersionCheck(svc *dynamodb.DynamoDB, tableName, id, newData string, expectedVersion int) error {

	updateExpression := "SET #data = :newData, #version = :newVersion"
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":newData": {
			S: aws.String(newData),
		},
		":newVersion": {
			N: aws.String(fmt.Sprintf("%d", expectedVersion+1)),
		},
		":expectedVersion": {
			N: aws.String(fmt.Sprintf("%d", expectedVersion)),
		},
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ConditionExpression:       aws.String("#version = :expectedVersion"),
		ExpressionAttributeNames: map[string]*string{
			"#data":    aws.String("data"),
			"#version": aws.String("version"),
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "ConditionalCheckFailedException" {
				return fmt.Errorf("Error: The current version is different from the expected version")
			}
		}
		return err
	}

	return nil
}
