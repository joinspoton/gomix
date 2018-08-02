package parse

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// BodyType - Type of the request body
type BodyType struct{}

// LambdaRequestBody - Parse a lambda http POST request body
func LambdaRequestBody(request events.APIGatewayProxyRequest) (BodyType, error) {
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	var parsedBody BodyType
	for {
		if err := decoder.Decode(&parsedBody); err == io.EOF {
			break
		} else if err != nil {
			return parsedBody, err
		}
	}
	return parsedBody, nil
}
