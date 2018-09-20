package mixpanel

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/joinspoton/gomix/utilities/paramstore"
)

// TimeInterval - Date range - YYYY-MM-DD
type TimeInterval struct {
	Start string
	End   string
}

// JQLQuery - parameters passed into jql query
type JQLQuery struct {
	Interval TimeInterval
	Event    string
	GroupBy  string
	OrderBy  string
	Filter   string
}

var (
	// mpurl - Mixpanel Query URL
	mpurl = "https://mixpanel.com/api/2.0/jql"
	// ErrorParameterStore - Parameter Store Error
	ErrorParameterStore = errors.New("Unable to get mixpanel secret from AWS paramater store")
)

// QueryMixpanel - hit mixpanel with jql query and return response
func QueryMixpanel(query JQLQuery) (string, error) {
	mpSecret, err := paramstore.GetConfig("/production/mixpanel/secret/b64")
	if err != nil {
		return "", ErrorParameterStore
	}
	dateRange := "{\"from_date\": \"" + query.Interval.Start + "\", \"to_date\": \"" + query.Interval.End + "\"}"
	jqlQuery := `
    function main() {
      return Events(params)
        .filter((event) => (
          event.name === '` + query.Event + `' ` + query.Filter + `
        ))
        .groupBy([` + query.GroupBy + `],mixpanel.reducer.count())
        .sortDesc('value')
    }
	`
	fmt.Println("THE JQL QUERY -------")
	fmt.Println(jqlQuery)

	data := url.Values{}
	data.Set("params", dateRange)
	data.Add("script", jqlQuery)
	req, err := http.NewRequest("POST", mpurl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Basic "+mpSecret)
	client := &http.Client{}
	var response string
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

//
func RawJQLQuery(jql string) (string, error) {
	return nil, nil
}
