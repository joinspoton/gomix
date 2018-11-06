package paramstore

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetConfig - Get config from paramstore
func GetConfig(key string) (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-west-2")},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	keyname := key
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}
	value := *param.Parameter.Value
	return value, nil
}

// GetJSONConfig - Get a map[string]interface{} object from paramstore
func GetJSONConfig(key string) (map[string]interface{}, error) {
	jsonString, _ := GetConfig(key)
	var config map[string]interface{}
	json.Unmarshal([]byte(jsonString), &config)
	return config, nil
}

// GetJSONArrayConfig - Get a []map[string]string object from paramstore
func GetJSONArrayConfig(key string) ([]map[string]string, error) {
	jsonString, _ := GetConfig(key)
	var config []map[string]string
	json.Unmarshal([]byte(jsonString), &config)
	return config, nil
}
