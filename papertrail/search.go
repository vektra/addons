package papertrail

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/vektra/cypress"
)

const cAPIRoot = "https://papertrailapp.com/api/v1/events/search.json"

type APIClient struct {
	*http.Client
	Token       string
	Options     *EventsOptions
	EventBuffer chan *Event
}

type EventsOptions struct {
	Q        string `url:"q,omitempty"`
	GroupId  string `url:"group_id,omitempty"`
	SystemId string `url:"system_id,omitempty"`
	MinId    string `url:"min_id,omitempty"`
	MaxId    string `url:"max_id,omitempty"`
	MinTime  string `url:"min_time,omitempty"`
	MaxTime  string `url:"max_time,omitempty"`
	Tail     string `url:"tail,omitempty"`
}

type EventsResponse struct {
	Events           *[]Event `json:"events"`
	MinId            string   `json:"min_id"`
	MaxId            string   `json:"max_id"`
	ReachedBeginning bool     `json:"reached_beginning"`
	ReachedTimeLimit bool     `json:"reached_time_limit"`
}

type Event struct {
	Id                string `json:"id'`
	ReceivedAt        string `json:"received_at"`
	DisplayReceivedAt string `json:"display_received_at"`
	SourceName        string `json:"source_name"`
	SourceId          uint32 `json:"source_id"`
	SourceIp          string `json:"source_ip"`
	Facility          string `json:"facility"`
	Severity          string `json:"severity"`
	Program           string `json:"program"`
	Message           string `json:"message"`
}

func NewAPIClient(token string, options *EventsOptions, bufferSize int) *APIClient {
	return &APIClient{
		Client:      &http.Client{},
		Token:       token,
		Options:     options,
		EventBuffer: make(chan *Event, bufferSize),
	}
}

func (api *APIClient) Search(o *EventsOptions) (*[]Event, error) {
	url := cAPIRoot

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

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Papertrail-Token", api.Token)

	resp, err := api.Do(req)
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
	if err != nil {
		return nil, err
	}

	return events.Events, nil
}

func (api *APIClient) Generate() (*cypress.Message, error) {
	select {

	case event := <-api.EventBuffer:

		var message cypress.Message

		err := json.Unmarshal([]byte(event.Message), &message)
		if err != nil {
			return nil, err
		}

		return &message, nil

	default:
		events, err := api.Search(&EventsOptions{})
		if err != nil {
			return nil, err
		}

		for _, event := range *events {
			api.EventBuffer <- &event
		}

		return api.Generate()
	}
}

func (api *APIClient) Close() error {
	// need to do anything here?
	return nil
}
