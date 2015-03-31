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
	*http.Client
	RootURL     string
	Options     *EventsOptions
	EventBuffer chan *cypress.Message
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

func NewAPIClient(key, host, log string, options *EventsOptions, bufferSize int) (*APIClient, error) {
	root := fmt.Sprintf("%s/%s/hosts/%s/%s/", cAPIRoot, key, host, log)

	url, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		Client:      &http.Client{},
		RootURL:     url.String(),
		Options:     options,
		EventBuffer: make(chan *cypress.Message, bufferSize),
	}, nil
}

func (api *APIClient) SetDefaultOptions(o *EventsOptions) *EventsOptions {
	if o.Start == 0 {
		o.Start = api.Options.Start
	}
	if o.End == 0 {
		o.End = api.Options.End
	}
	if o.Filter == "" {
		o.Filter = api.Options.Filter
	}
	if o.Limit == 0 {
		o.Limit = api.Options.Limit
	}

	return o
}

func (api *APIClient) EncodeURL(o *EventsOptions) string {
	url := api.RootURL

	v, _ := query.Values(o)
	if q := v.Encode(); q != "" {
		url = url + "?" + q
	}

	return url
}

func (api *APIClient) GetBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func NewEvents(body []byte) ([]*cypress.Message, error) {
	var events EventsResponse
	err := json.Unmarshal(body, &events)

	if err == nil {
		if events.Response == "error" {
			message := fmt.Sprintf("Logentries error: %s", events.Response, events.Reason)
			return nil, errors.New(message)
		} else {
			// Ok but no events
			return nil, errors.New("Logentires error: No events")
		}

	} else {
		// Log lines sent back verbatim, not proper JSON
		logs := bytes.Split(body, []byte("\n"))

		var events []*cypress.Message

		for _, log := range logs {
			if string(log) != "" {
				var message cypress.Message

				err = json.Unmarshal(log, &message)
				if err != nil {
					message = *cypress.Log()
					message.AddString("message", string(log))
				}

				events = append(events, &message)
			}
		}

		return events, nil
	}

	return nil, nil
}

func (api *APIClient) Search(o *EventsOptions) ([]*cypress.Message, error) {
	opts := api.SetDefaultOptions(o)
	url := api.EncodeURL(opts)

	body, err := api.GetBody(url)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	return NewEvents(body)
}

func milliseconds(t time.Time) int {
	nanos := t.UnixNano()
	millis := nanos / 1000000
	return int(millis)
}

func (api *APIClient) BufferEvents(events []*cypress.Message) error {
	for _, event := range events {
		select {

		case api.EventBuffer <- event:
			api.Options.Start = milliseconds(event.GetTimestamp().Time())

		default:
			break
		}
	}

	return nil
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

		api.BufferEvents(events)

		return api.Generate()
	}
}

func (api *APIClient) Close() error {
	close(api.EventBuffer)
	return nil
}
