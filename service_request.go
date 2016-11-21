package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type serviceRequest struct {
	ServiceRequestID  string  `json:"service_request_id"`
	Status            string  `json:"status"`
	ServiceName       string  `json:"service_name"`
	ServiceCode       string  `json:"service_code"`
	AgencyResponsible string  `json:"agency_responsible"`
	Description       string  `json:"description"`
	RequestedDatetime string  `json:"requested_datetime"`
	UpdatedDatetime   string  `json:"updated_datetime"`
	Address           string  `json:"address"`
	Lat               float64 `json:"lat"`
	Long              float64 `json:"long"`
	MediaURL          string  `json:"media_url"`
}

func latestRequests(reqURL string, startDate time.Time) ([]serviceRequest, error) {
	format := "2006-01-02T15:04:05-07:00"
	timeStr := startDate.Format(format)
	v := url.Values{}
	v.Set("start_date", timeStr)

	log.Printf("checking time %s", timeStr)

	resp, err := http.Get(reqURL + "?" + v.Encode())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var requests []serviceRequest

	err = json.NewDecoder(resp.Body).Decode(&requests)

	if err != nil {
		return nil, err
	}

	return requests, nil
}
