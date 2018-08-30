package paramstore

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetConfig - Get config from paramstore
func GetConfig(key string) ([]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-west-2")},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	keyname := key
	withDecryption := false
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	value := *param.Parameter.Value
	config := strings.Split(value, ",")
	return config, nil
}
