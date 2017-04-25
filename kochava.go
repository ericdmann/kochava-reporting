package kochava

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	reportingVersion  = "v1.3"
	pathSummary       = "summary"
	pathDetail        = "detail"
	pathStatus        = "progress"
	reportingEndpoint = "https://reporting.api.kochava.com"
)

type ReportingRequest struct {
	Token           string   `json:"token,omitempty"`
	Type            string   `json:"-"`
	APIKey          string   `json:"api_key,omitempty"`
	AppGUID         string   `json:"app_guid,omitempty"`
	Report          string   `json:"report,omitempty"`
	TimeStart       string   `json:"time_start,omitempty"`
	TimeEnd         string   `json:"time_end,omitempty"`
	Traffic         []string `json:"traffic,omitempty"`
	TimeZone        string   `json:"time_zone,omitempty"`
	DeliveryMethod  []string `json:"delivery_method,omitempty"`
	DeliveryFormat  string   `json:"delivery_format,omitempty"`
	TrafficGrouping []string `json:"traffic_grouping,omitempty"`
	Timeseries      string   `json:"time_series,omitempty"`
	Notify          []string `json:"notify,omitempty"`
	Event           string   `json:"event,omitempty"`
}

type Report struct {
}

type ReportResponse struct {
	Status          string           `json:"status, omitempty"`
	ReportLink      string           `json:"report, omitempty"`
	StatusDate      string           `json:"status_date,omitempty"`
	Progress        string           `json:"progress,omitempty"`
	ReportTemplates []ReportTemplate `json:"template_values,omitempty"`
	ReportToken     string           `json:"report_token"`
	Error           string           `json:"error"`
}

type ReportTemplate struct {
	ReportType       string   `json:"report_type"`
	ReportComponent  string   `json:"report_section"`
	DefaultColumns   []string `json:"columns_selected"`
	AvailableColumns []string `json:"available_columns"`
}
type Kochava struct {
	HClient          http.Client
	APIKey           string
	Version          string
	QueuedReports    []Report
	AppGUID          string
	CompletedReports []Report
	ErroredReports   []Report
	Templates        []ReportTemplate
}

func NewClient(apiKey string, appGUID string) (Kochava, error) {
	var err error

	client := Kochava{}
	client.HClient = http.Client{}
	client.AppGUID = appGUID
	client.APIKey = apiKey
	client.Templates, err = client.RetrieveTemplates()
	return client, err
}

func (k *Kochava) RetrieveTemplates() ([]ReportTemplate, error) {
	var repResponse ReportResponse

	request := ReportingRequest{}
	request.APIKey = k.APIKey
	request.AppGUID = k.AppGUID

	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", reportingEndpoint+"/"+reportingVersion+"/reportcolumns", strings.NewReader(string(body)))
	resp, err := k.HClient.Do(req)

	if err != nil {
		return repResponse.ReportTemplates, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		err = errors.New("Unable to fetch columns")
	}

	json.Unmarshal(respBody, &repResponse)

	return repResponse.ReportTemplates, err

}
