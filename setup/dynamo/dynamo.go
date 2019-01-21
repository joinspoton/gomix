package dynamo

import (
	"fmt"
	"gomix/utilities/system"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func getClient() *dynamodb.DynamoDB {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create DynamoDB client
	return dynamodb.New(sess)
}

// CreateItems - insert items into a DynamoDB table
func CreateItems(items []interface{}, table string) {
	// TODO: Autogenerate primary key
	svc := getClient()

	itemsLength := len(items)
	var wg sync.WaitGroup
	wg.Add(itemsLength)

	for i := 0; i < itemsLength; i++ {
		go func(i int) {
			defer wg.Done()
			item := items[i]

			av, err := dynamodbattribute.MarshalMap(item)
			id, _ := dynamodbattribute.Marshal(system.CreateUUID())
			av["id"] = id

			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String(table),
			}
			_, err = svc.PutItem(input)

			if err != nil {
				fmt.Println("Got error calling PutItem:")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("Successfully added %+v to %s\n", item, table)
		}(i)
	}

	wg.Wait()
}

// GetAllItems - Retrive every items in a table
func GetAllItems(table string) []map[string]interface{} {
	svc := getClient()

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}

	items := []map[string]interface{}{}
	svc.ScanPages(params,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			for _, unmarshalItem := range page.Items {
				item := make(map[string]interface{})
				dynamodbattribute.UnmarshalMap(unmarshalItem, &item)
				items = append(items, item)
			}
			return !lastPage
		})

	return items
}

// RemoveTable - Delete the table from dynamo (DeleteTable is a dynamo function)
func RemoveTable(table string) error {
	svc := getClient()

	// Build the required object
	var tableObj dynamodb.DeleteTableInput
	tableObj.TableName = &table

	svc.DeleteTable(&tableObj)

	return nil
}

// CreateTable - Create the table in dynamo with primary key
func CreateTable(table string, primaryKeyName string, primaryKeyType string, readVal int64, writeVal int64) error {
	var err error
	err = nil
	svc := getClient()

	// Build the keys
	var primaryKeyAttribute dynamodb.AttributeDefinition
	var attributes []*dynamodb.AttributeDefinition
	primaryKeyAttribute.AttributeName = &primaryKeyName
	primaryKeyAttribute.AttributeType = &primaryKeyType
	attributes = append(attributes, &primaryKeyAttribute)

	// Indicate the primary key. HASH indicates it is a single primary key (RANGE for composite keys)
	var test dynamodb.KeySchemaElement
	data := "HASH"
	test.AttributeName = &primaryKeyName
	test.KeyType = &data
	var keySchema []*dynamodb.KeySchemaElement
	keySchema = append(keySchema, &test)

	// Indicate throughputs
	var throughput dynamodb.ProvisionedThroughput
	throughput.ReadCapacityUnits = &readVal
	throughput.WriteCapacityUnits = &writeVal

	// Build the required object
	var tableObj dynamodb.CreateTableInput
	tableObj.TableName = &table
	tableObj.AttributeDefinitions = attributes
	tableObj.KeySchema = keySchema
	tableObj.ProvisionedThroughput = &throughput

	_, err = svc.CreateTable(&tableObj)
	return err
}
