package influx

import (
	"errors"
	"fmt"
	"gomix/utilities/paramstore"
	"gomix/utilities/system"

	"github.com/influxdata/influxdb/client/v2"
)

// ConfigType - Required to connect to influx db
type ConfigType struct {
	host     string
	username string
	password string
}

// ErrInfluxConnection - Error thrown when connection failes
var ErrInfluxConnection = errors.New("INFLUX_FAILED_TO_CONNECT")

// Connect - Create a new InfluxDB connection without credentials to ec2-34-209-159-101.us-west-2.compute.amazonaws.com
func Connect() (clnt client.Client, err error) {
	path := fmt.Sprintf("/%s/influx/", system.GetEnv("stage", "staging"))

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		fmt.Sprintf("http://%s:%s", host, port),
		username,
		password,
	)
}

// ConnectToProd - Create a new InfluxDB connection using production credentials
func ConnectToProd() (clnt client.Client, err error) {
	path := "/production/influx/"

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		fmt.Sprintf("http://%s:%s", host, port),
		username,
		password,
	)
}

// Connect6c4b - Create a new InfluxDB connection without credentials to ec2-34-216-111-45.us-west-2.compute.amazonaws.com
func Connect6c4b() (clnt client.Client, err error) {
	path := fmt.Sprintf("/%s/influx/6c4b/", system.GetEnv("stage", "staging"))

	host, _ := paramstore.GetConfig(path + "host")
	port, _ := paramstore.GetConfig(path + "port")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		fmt.Sprintf("http://%s:%s", host, port),
		username,
		password,
	)
}

// ManuallyConnect - Create a new InfluxDB connection with credentials
func ManuallyConnect(host string, username string, password string) (clnt client.Client, err error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     host,
		Username: username,
		Password: password,
	})
	if err != nil {
		fmt.Println("Unable to connect to influxDB:")
		fmt.Println("Host: ", host)
		fmt.Println("Username: ", username)
		fmt.Println("Password: ", password)
		return nil, ErrInfluxConnection
	}
	return c, nil
}
