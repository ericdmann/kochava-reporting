package kochava

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (k *Kochava) CheckRequest(token string) (ReportResponse, error) {
	var err error
	var repResponse ReportResponse
	var repReq ReportingRequest

	repReq.APIKey = k.APIKey
	repReq.AppGUID = k.AppGUID
	repReq.Token = token
	body, err := json.Marshal(repReq)

	if err != nil {
		return repResponse, err
	}

	req, err := http.NewRequest("POST", reportingEndpoint+"/"+reportingVersion+"/"+pathStatus, strings.NewReader(string(body)))
	resp, err := k.HClient.Do(req)

	if err != nil {
		return repResponse, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		err = errors.New("Unable to submit status request: " + resp.Status + string(respBody))
		return repResponse, err
	}

	json.Unmarshal(respBody, &repResponse)

	if repResponse.Status == "Error" {
		return repResponse, errors.New(repResponse.Error)
	}
	return repResponse, err
}

func (k *Kochava) NewDetailRequest(traffic string,
	deliveryFormat string,
	notify []string,
	timeStart time.Time,
	timeEnd time.Time,
	timezone string) (ReportingRequest, error) {

	var err error
	var req ReportingRequest
	req.APIKey = k.APIKey
	req.AppGUID = k.AppGUID
	req.DeliveryFormat = deliveryFormat
	req.DeliveryMethod = []string{"S3link"}
	req.Notify = notify
	req.TimeStart = strconv.Itoa(int(timeStart.Unix()))
	req.TimeEnd = strconv.Itoa(int(timeEnd.Unix()))
	req.TimeZone = timezone
	req.Traffic = []string{traffic}
	req.Type = "detail"

	return req, err
}

func (k *Kochava) NewSummaryRequest(traffic string,
	timeseries string,
	trafficGrouping []string,
	deliveryFormat string,
	notify []string,
	timeStart time.Time,
	timeEnd time.Time,
	timezone string) (ReportingRequest, error) {

	var err error
	var req ReportingRequest
	req.APIKey = k.APIKey
	req.AppGUID = k.AppGUID
	req.DeliveryFormat = deliveryFormat
	req.DeliveryMethod = []string{"S3link"}
	req.Notify = notify
	req.TimeStart = strconv.Itoa(int(timeStart.Unix()))
	req.TimeEnd = strconv.Itoa(int(timeEnd.Unix()))
	req.TimeZone = timezone
	req.Timeseries = timeseries
	req.TrafficGrouping = trafficGrouping
	req.Traffic = []string{traffic}
	req.Type = "summary"

	return req, err
}

func (k *Kochava) SubmitRequest(repReq ReportingRequest) (string, error) {
	var err error

	endpoint := pathSummary
	if repReq.Type == "detail" {
		endpoint = pathDetail
	}

	body, err := json.Marshal(repReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", reportingEndpoint+"/"+reportingVersion+"/"+endpoint, strings.NewReader(string(body)))
	resp, err := k.HClient.Do(req)

	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		err = errors.New("Unable to submit report: " + resp.Status + string(respBody))
		return "", err
	}

	var repResponse ReportResponse
	json.Unmarshal(respBody, &repResponse)

	if repResponse.Status == "Error" {
		return "", errors.New(repResponse.Error)
	}
	return repResponse.ReportToken, err
}
