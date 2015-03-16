package logentries

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/vektra/cypress"
)

const cAPIRoot = "https://pull.logentries.com/"

type APIClient struct {
	accountKey string
	host       string
	log        string
}

type EventsOptions struct {
	Start  int    `url:"start,omitempty"`
	End    int    `url:"end,omitempty"`
	Filter string `url:"filter,omitempty"`
	Limit  int    `url:"limit,omitempty"`
}

type EventsResponse struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
	Events   *[]cypress.Message
}

func (api *APIClient) apiRoot() string {
	root := cAPIRoot + api.accountKey + "/" + api.logAddr() + "/"

	url, err := url.Parse(root)
	if err != nil {
		panic(err)
	}

	return url.String()
}

func (api *APIClient) logAddr() string {
	return "hosts/" + api.host + "/" + api.log
}

func (api *APIClient) Search(o *EventsOptions) (*EventsResponse, error) {
	events, err := api.SearchEvents(o)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (api *APIClient) SearchEvents(o *EventsOptions) (*EventsResponse, error) {
	v, _ := query.Values(o)

	url := api.apiRoot()
	if query := v.Encode(); query != "" {
		url = url + "?" + query
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events EventsResponse

	if body[0] == []byte("{")[0] {
		err = json.Unmarshal(body, &events)
		if err != nil {

			if err.Error() == "invalid character '{' after top-level value" {

				logs := bytes.Split(body, []byte("\n"))

				for i := 0; i < len(logs); i++ {
					log := logs[i]
					var message cypress.Message

					json.Unmarshal(log, &message)
					if err != nil {
						return nil, err
					}

					newEvents := append(*events.Events, message)
					events.Events = &newEvents
				}

			} else {
				return nil, err
			}
		}
	}

	return &events, nil
}
