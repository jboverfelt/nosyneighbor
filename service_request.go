package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type serviceRequest struct {
	Address     string `json:"address"`
	Agency      string `json:"agency"`
	CreatedDate string `json:"createddate"`
	DueDate     string `json:"duedate"`
	Geolocation struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geolocation"`
	Latitude         string `json:"latitude"`
	Longitude        string `json:"longitude"`
	Method           string `json:"methodreceived"`
	Neighborhood     string `json:"neighborhood"`
	CouncilDistrict  string `json:"councildistrict"`
	PoliceDistrict   string `json:"policedistrict"`
	ServiceRequestID string `json:"servicerequestnum"`
	SRRecordID       string `json:"srrecordid"`
	Status           string `json:"srstatus"`
	ServiceName      string `json:"srtype"`
	StatusDate       string `json:"statusdate"`
	Zipcode          string `json:"zipcode"`
}

func latestRequests(reqURL string, startDate time.Time) ([]serviceRequest, error) {
	format := "2006-01-02T15:04:05.999"
	timeStr := startDate.Format(format)

	v := url.Values{}
	v.Set("$where", fmt.Sprintf("statusdate > '%s'", timeStr))
	v.Set("$order", "statusdate DESC")

	log.Printf("checking time %s", timeStr)

	resp, err := http.Get(reqURL + "?" + v.Encode())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var buf strings.Builder
	tee := io.TeeReader(resp.Body, &buf)

	var requests []serviceRequest
	err = json.NewDecoder(tee).Decode(&requests)

	if err != nil {
		log.Printf("json response: %s", buf.String())
		return nil, err
	}

	return requests, nil
}
