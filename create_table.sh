#!/bin/bash

table_name="TableExample"

export AWS_DEFAULT_REGION=us-east-1

aws dynamodb create-table --region us-east-1 --endpoint-url=http://localhost:4566 \
  --table-name $table_name \
  --attribute-definitions \
    AttributeName=id,AttributeType=S \
  --key-schema \
    AttributeName=id,KeyType=HASH \
  --provisioned-throughput \
    ReadCapacityUnits=1,WriteCapacityUnits=1

if [ $? -eq 0 ]; then
  echo "Table $table_name created successfully."
else
  echo "Error creating table $table_name."
fi
