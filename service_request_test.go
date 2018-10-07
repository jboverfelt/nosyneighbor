package main

import (
	"encoding/json"
	"strings"
	"testing"
)

var sample = `
[{"address":"1421 S CATON AVE, Baltimore City, 21227","agency":"Liquor Board","councildistrict":"10","createddate":"2018-10-07T05:50:59.000","duedate":"2018-11-24T05:50:59.000","geolocation":{"type":"Point","coordinates":[-76.665096547748,39.262600186508]},"latitude":"39.262600186508195","longitude":"-76.665096547747770","methodreceived":"API","neighborhood":"Morrell Park","policedistrict":"Southwestern","servicerequestnum":"18-00762026","srrecordid":"5004100000gKG77AAG","srstatus":"New","srtype":"BCLB-Liquor License Complaint","statusdate":"2018-10-07T05:50:59.000","zipcode":"21227"}
,{"address":"2633 BERYL AVE, Baltimore City, 21205","agency":"Housing","councildistrict":"13","createddate":"2018-10-07T05:44:52.000","duedate":"2018-10-22T05:44:52.000","geolocation":{"type":"Point","coordinates":[-76.580187818834,39.30293457882]},"latitude":"39.302934578820100","longitude":"-76.580187818834080","methodreceived":"API","neighborhood":"Biddle Street","policedistrict":"Eastern","servicerequestnum":"18-00762025","srrecordid":"5004100000gKG72AAG","srstatus":"New","srtype":"HCD-Vacant Building","statusdate":"2018-10-07T05:44:52.000","zipcode":"21205"}]
`

func TestUnmarshal(t *testing.T) {
	var requests []serviceRequest
	err := json.NewDecoder(strings.NewReader(sample)).Decode(&requests)

	if err != nil {
		t.Fatal(err)
	}
}
