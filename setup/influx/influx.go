package gomix

import (
	"errors"

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
func Connect(influxConfig InfluxConfigType) (clnt client.Client, err error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxConfig.host,
		Username: influxConfig.username,
		Password: influxConfig.password,
	})
	if err != nil {
		return nil, ErrInfluxConnection
	}
	return c, nil
}
