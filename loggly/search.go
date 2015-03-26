package loggly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/vektra/cypress"
)

const cAPIRootRSID = "loggly.com/apiv2/search"
const cAPIRootEvents = "loggly.com/apiv2/events"

type APIClient struct {
	*http.Client
	Username      string
	Password      string
	RSIDRootURL   string
	EventsRootURL string
	RSIDOptions   *RSIDOptions
	EventsOptions *EventsOptions
	EventBuffer   chan *Event
}

func NewAPIClient(account, username, password string, ro *RSIDOptions, eo *EventsOptions, bufferSize int) (*APIClient, error) {
	rsid := fmt.Sprintf("http://%s.%s", account, cAPIRootRSID)
	rsidUrl, err := url.Parse(rsid)
	if err != nil {
		return nil, err
	}

	events := fmt.Sprintf("http://%s.%s", account, cAPIRootEvents)
	eventsUrl, err := url.Parse(events)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		Client:        &http.Client{},
		Username:      username,
		Password:      password,
		RSIDRootURL:   rsidUrl.String(),
		EventsRootURL: eventsUrl.String(),
		RSIDOptions:   ro,
		EventsOptions: eo,
		EventBuffer:   make(chan *Event, bufferSize),
	}, nil
}

type RSIDOptions struct {
	Q     string `url:"q,omitempty"`
	From  string `url:"from,omitempty"`
	Until string `url:"until,omitempty"`
	Order string `url:"order,omitempty"`
	Size  uint   `url:"size,omitempty"`
}

type RSIDResponse struct {
	RSID `json:"rsid"`
}

type RSID struct {
	Status      string  `json:"status"`
	DateFrom    uint    `json:"date_from"`
	ElapsedTime float64 `json:"elapsed_time"`
	DateTo      uint    `json:"date_to"`
	ID          string  `json:"id"`
}

type EventsOptions struct {
	RSID    string `url:"rsid"`
	Page    uint   `url:"page,omitempty"`
	Format  string `url:"format,omitempty"`
	Columns string `url:"columns,omitempty"`
}

type EventsResponse struct {
	TotalEvents uint     `json:"total_events"`
	Page        uint     `json:"page"`
	Events      []*Event `json:"events"`
}

type Event struct {
	Timestamp uint   `json:"timestamp"`
	Logmsg    string `json:"logmsg"`
}

func (api *APIClient) Search(ro *RSIDOptions, eo *EventsOptions) ([]*Event, error) {
	rsid, err := api.SearchRSID(ro)
	if err != nil {
		return nil, err
	}

	eo.RSID = rsid.ID

	events, err := api.SearchEvents(eo)
	if err != nil {
		return nil, err
	}

	return events.Events, nil
}

func (api *APIClient) SearchRSID(o *RSIDOptions) (*RSIDResponse, error) {
	url := api.RSIDRootURL

	v, _ := query.Values(o)
	if q := v.Encode(); q != "" {
		url = url + "?" + q

	} else {
		// Use default options
		v, _ = query.Values(api.RSIDOptions)
		if q = v.Encode(); q != "" {
			url = url + "?" + q
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(api.Username, api.Password)

	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rsid RSIDResponse
	err = json.Unmarshal(body, &rsid)
	if err != nil {
		return nil, err
	}

	return &rsid, nil
}

func (api *APIClient) SearchEvents(o *EventsOptions) (*EventsResponse, error) {
	url := api.EventsRootURL

	v, _ := query.Values(o)
	if q := v.Encode(); q != "" {
		url = url + "?" + q

	} else {
		// Use default options
		v, _ = query.Values(api.EventsOptions)
		if q = v.Encode(); q != "" {
			url = url + "?" + q
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(api.Username, api.Password)

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

	return &events, nil
}

func (api *APIClient) Generate() (*cypress.Message, error) {
	select {

	case event := <-api.EventBuffer:
		var message *cypress.Message
		err := json.Unmarshal([]byte(event.Logmsg), message)
		if err != nil {
			message = cypress.Log()
			message.Add("message", event.Logmsg)
		}

		return message, nil

	case <-time.After(time.Second * 1):
		return nil, nil

	default:
		events, err := api.Search(api.RSIDOptions, api.EventsOptions)
		if err != nil {
			return nil, err
		}

		for _, event := range events {
			select {

			case api.EventBuffer <- event:
				api.EventsOptions.Page = api.EventsOptions.Page + 1

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
