package paramstore

import (
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

	fmt.Println("key", key)
	c, _ := sess.Config.Credentials.Get()
	fmt.Println("sess.Config.Credentials.Get()", c)

	// fmt.Println("Trying to request...")
	// resp, err := http.Get("https://github.com/")
	// if err != nil {
	// 	fmt.Println("Error!")
	// }
	// fmt.Println("After requesting...")
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println("body", string(body))

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	fmt.Printf("ssmsvc.ClientInfo: \n%+v\n", ssmsvc.ClientInfo)
	d, _ := ssmsvc.Config.Credentials.Get()
	fmt.Println("ssmsvc.Config.Credentials.Get()", d)
	keyname := key
	withDecryption := false
	fmt.Println("Before ssmsvc.GetParameter()")
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	fmt.Println("After ssmsvc.GetParameter()")
	fmt.Println("param.GoString()", param.GoString())
	if err != nil {
		fmt.Println("err", err)
	}

	value := *param.Parameter.Value
	fmt.Println("value", value)
	return value, nil
}
