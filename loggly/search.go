package loggly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const cAPIRootRSID = "loggly.com/apiv2/search"
const cAPIRootEvents = "loggly.com/apiv2/events"

type APIClient struct {
	account  string
	username string
	password string
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
	totalEvents int
	page        int
	events      *[]Event
}

type Event struct {
	tags      *[]string
	timestamp string
	logmsg    string
	logTypes  *[]string
	id        string
	// EventDetail?
}

func (api *APIClient) apiRootRSID() string {
	root := "http://" + api.account + "." + cAPIRootRSID

	url, err := url.Parse(root)
	if err != nil {
		panic(err)
	}

	return url.String()
}

func (api *APIClient) apiRootEvents() string {
	root := "http://" + api.account + "." + cAPIRootEvents

	url, err := url.Parse(root)
	if err != nil {
		panic(err)
	}

	return url.String()
}

func (api *APIClient) Search(ro *RSIDOptions, eo *EventsOptions) (*EventsResponse, error) {
	rsid, err := api.SearchRSID(ro)
	if err != nil {
		return nil, err
	}

	eo.RSID = rsid.ID

	events, err := api.SearchEvents(eo)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (api *APIClient) SearchRSID(o *RSIDOptions) (*RSIDResponse, error) {
	client := &http.Client{}

	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	url := api.apiRootRSID()
	if q := v.Encode(); q != "" {
		url = url + "?" + q
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(api.username, api.password)

	resp, err := client.Do(req)
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
	client := &http.Client{}

	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}

	url := api.apiRootEvents()
	if q := v.Encode(); q != "" {
		url = url + "?" + q
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(api.username, api.password)

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
