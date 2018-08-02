package parse

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// BodyType - Type of the request body
type BodyType struct{}

// ErrParseRequestBody - error when parsing request body
var ErrParseRequestBody = errors.New("ERROR_PARSE_REQUEST_BODY")

// LambdaRequestBody - Parse a lambda http POST request body
func LambdaRequestBody(request events.APIGatewayProxyRequest) (BodyType, error) {
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	var parsedBody BodyType
	for {
		if err := decoder.Decode(&parsedBody); err == io.EOF {
			break
		} else if err != nil {
			return parsedBody, ErrParseRequestBody
		}
	}
	return parsedBody, nil
}
