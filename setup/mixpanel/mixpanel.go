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

// RawJQLQuery - Takes any parameters and JQL, pings MixPanel, and returns the result
func RawJQLQuery(params string, script string) (string, error) {
	mpSecret, err := paramstore.GetConfig("/production/mixpanel/secret/b64")
	if err != nil {
		return "", ErrorParameterStore
	}

	data := url.Values{}
	if params != "" {
		data.Set("params", params)
	}
	data.Set("script", script)

	req, _ := http.NewRequest("POST", mpurl, strings.NewReader(data.Encode()))
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

func ReadMixpanelData(body []byte) []byte {
	var openCount = 0
	var newBody []byte
	var flag = 0
	var err error
	for index, element := range body[:len(body)-3] {
		if element == 123 {
			openCount++
		}
		if openCount < 3 {
			newBody = append(newBody, element)
		}
		if openCount == 3 {
			newBody = append(newBody, byte(91))
			openCount++
		}
		if openCount >= 3 && element == 34 {
			//subsqeunt open brackets
			if flag%3 == 0 {
				newBody = append(newBody, byte(123))
				var quoteCount = 0
				var nameStart = index
				var nameEnd = index
				var dateStart = index
				var dateEnd = index
				var countStart = index
				var countEnd = index
				var site []byte
				var dte []byte
				var count []byte
				var appendList []byte
				var siteStr []byte
				var dateStr []byte
				var countStr []byte
				for subIndex, subElement := range body[index:] {
					subIndex += index
					if subElement == 34 {
						quoteCount++
					}
					if quoteCount == 1 && subElement == 34 {
						nameStart = subIndex
					} else if quoteCount == 2 && subElement == 34 {
						nameEnd = subIndex
						site = body[nameStart+1 : nameEnd]
					} else if quoteCount == 3 && subElement == 34 {
						dateStart = subIndex
					} else if quoteCount == 4 && subElement == 34 {
						dateEnd = subIndex
						dte = body[dateStart+1 : dateEnd]
					} else if quoteCount == 4 && subElement == 58 {
						countStart = subIndex
					} else if quoteCount == 4 && subElement == 125 {
						countEnd = subIndex
						count = body[countStart+2 : countEnd]
						quoteCount++
						break
					}
					if quoteCount > 5 {
						break
					}
				}
				siteStr = append(siteStr, byte(34), byte(115), byte(105), byte(116), byte(101), byte(34), byte(58))
				appendList = append(appendList, siteStr...)
				appendList = append(appendList, byte(34))
				appendList = append(appendList, site...)
				appendList = append(appendList, byte(34), byte(44))

				dateStr = append(dateStr, byte(34), byte(100), byte(97), byte(116), byte(101), byte(34), byte(58))
				appendList = append(appendList, dateStr...)
				appendList = append(appendList, byte(34))
				appendList = append(appendList, dte...)
				appendList = append(appendList, byte(34), byte(44))

				countStr = append(countStr, byte(34), byte(99), byte(111), byte(117), byte(110), byte(116), byte(34), byte(58))
				appendList = append(appendList, countStr...)
				appendList = append(appendList, count...)
				newBody = append(newBody, appendList...)
				newBody = append(newBody, byte(125))
				if body[countEnd+2] != 125 {
					newBody = append(newBody, byte(44))
				}

			}
			flag++
			fmt.Println(flag)
		}

		if body[index] == 125 && body[index+1] == 125 && body[index+2] == 125 {
			newBody = append(newBody, byte(93), byte(125), byte(125))
		}

	}
	if err != nil {
		panic(err)
	}
	return newBody

}
