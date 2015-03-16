package papertrail

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

const cAPIRoot = "https://papertrailapp.com/api/v1/events/search.json"

type APIClient struct {
	token string
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

func (api *APIClient) Search(o *EventsOptions) (*EventsResponse, error) {
	events, err := api.SearchEvents(o)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (api *APIClient) SearchEvents(o *EventsOptions) (*EventsResponse, error) {
	client := &http.Client{}

	v, _ := query.Values(o)

	url := cAPIRoot
	if query := v.Encode(); query != "" {
		url = url + "?" + query
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Papertrail-Token", api.token)

	resp, err := client.Do(req)
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

	return &events, nil
}
