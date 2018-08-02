package influx

import (
	"errors"
	"fmt"

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

// Connect - chodie
func Connect(host string, username string, password string) (clnt client.Client, err error) {
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
		fmt.Println("Unable to connect to influxDB:")
		return nil, ErrInfluxConnection
	}
	return c, nil
}
