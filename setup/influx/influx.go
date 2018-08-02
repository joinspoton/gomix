package gomix

import (
	"errors"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/influxdata/influxdb/client/v2"
)

// InfluxConfigType - Required to connect to influx db
type InfluxConfigType struct {
	host     string
	username string
	password string
}

// ErrInfluxConnection - Error thrown when connection failes
var ErrInfluxConnection = errors.New("INFLUX_FAILED_TO_CONNECT")

// Connect - chodie
func Connect(influxConfig struct{}) (clnt client.Client, err error) {
	config := InfluxConfigType{}
	mergo.Merge(&config, influxConfig)
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.host,
		Username: config.username,
		Password: config.password,
	})
	if err != nil {
		fmt.Println("Unable to connect to influxDB:")
		fmt.Println("Host: ", config.host)
		fmt.Println("Username: ", config.username)
		fmt.Println("Password: ", config.password)
		fmt.Println("Unable to connect to influxDB:")
		return nil, ErrInfluxConnection
	}
	return c, nil
}
