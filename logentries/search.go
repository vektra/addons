package logentries

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/vektra/cypress"
)

const cAPIRoot = "https://pull.logentries.com"

type APIClient struct {
	RootURL     string
	Options     *EventsOptions
	EventBuffer chan *cypress.Message
}

func NewAPIClient(key, host, log string, options *EventsOptions, bufferSize int) (*APIClient, error) {
	root := fmt.Sprintf("%s/%s/hosts/%s/%s/", cAPIRoot, key, host, log)

	url, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		RootURL:     url.String(),
		Options:     options,
		EventBuffer: make(chan *cypress.Message, bufferSize),
	}, nil
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
	Events   []*cypress.Message
}

func (api *APIClient) Search(o *EventsOptions) ([]*cypress.Message, error) {
	url := api.RootURL

	v, _ := query.Values(o)
	if q := v.Encode(); q != "" {
		url = url + "?" + q

	} else {
		// Use default options
		v, _ = query.Values(api.Options)
		if q = v.Encode(); q != "" {
			url = url + "?" + q
		}
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
	err = json.Unmarshal(body, &events)

	if err == nil {
		if events.Response == "error" {
			message := fmt.Sprintf("Logentries error: %s", events.Response, events.Reason)
			return nil, errors.New(message)
		} else {
			// Ok but no events
			return nil, errors.New("Logentires error: No events")
		}

	} else if err.Error() == "invalid character '{' after top-level value" {
		// Log lines sent back verbatim, not proper JSON
		logs := bytes.Split(body, []byte("\n"))

		var events []*cypress.Message

		for i := 0; i < len(logs)-1; i++ {
			log := logs[i]
			var message *cypress.Message

			err = json.Unmarshal(log, message)
			if err != nil {
				message = cypress.Log()
				message.Add("message", log)
			}

			events = append(events, message)
		}

		return events, nil

	} else {
		// Unknown error
		return nil, err
	}
}

func milliseconds(t time.Time) int {
	nanos := t.UnixNano()
	millis := nanos / 1000000
	return int(millis)
}

func (api *APIClient) Generate() (*cypress.Message, error) {
	select {

	case event := <-api.EventBuffer:
		return event, nil

	case <-time.After(time.Second * 1):
		return nil, nil

	default:
		events, err := api.Search(api.Options)
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			select {

			case api.EventBuffer <- event:
				api.Options.Start = milliseconds(event.GetTimestamp().Time())

			case <-time.After(time.Second * 1):
				break
			}
		}

		return api.Generate()
	}
}

func (api *APIClient) Close() error {
	// need to do anything here?
	return nil
}
