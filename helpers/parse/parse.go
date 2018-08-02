package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// BodyType - Type of the request body
type BodyType struct {
	Site     string
	Calendar string
}

// ErrParseRequestBody - error when parsing request body
var ErrParseRequestBody = errors.New("ERROR_PARSE_REQUEST_BODY")

// LambdaRequestBody - Parse a lambda http POST request body
func LambdaRequestBody(request events.APIGatewayProxyRequest) (BodyType, error) {
	fmt.Println(request.Body)
	decoder := json.NewDecoder(strings.NewReader(request.Body))
	var parsedBody BodyType
	for {
		if err := decoder.Decode(&parsedBody); err == io.EOF {
			break
		} else if err != nil {
			return parsedBody, ErrParseRequestBody
		}
	}
	fmt.Println(reflect.TypeOf(parsedBody))
	return parsedBody, nil
}

// decoder := json.NewDecoder(strings.NewReader(request.Body))
// var loadPoint LoadPoint
// for {
// 	if err := decoder.Decode(&loadPoint); err == io.EOF {
// 		break
// 	} else if err != nil {
// 		return LoadPoint{}, err
// 	}
// }
// return loadPoint, nil
