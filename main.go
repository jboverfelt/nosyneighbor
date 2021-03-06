package main

import (
	"fmt"
	"log"
	"os"

	"googlemaps.github.io/maps"

	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type env struct {
	requestURL      string
	fromAddr        string
	toAddr          string
	dist            distancer
	mail            emailer
	home            string
	councilDistrict string
	threshold       meters
	checkInterval   time.Duration
}

func main() {
	curEnv := setupEnv()

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	for {
		startDate := time.Now().Add(-curEnv.checkInterval).In(loc)
		err := checkLatestRequests(curEnv, startDate)

		if err != nil {
			curEnv.mail.SendErrorEmail(curEnv.fromAddr, curEnv.toAddr, err)
			return
		}

		time.Sleep(curEnv.checkInterval)
	}

}

func setupEnv() *env {
	mapsKey := os.Getenv("NOSY_MAPS_KEY")
	mailgunPubKey := os.Getenv("NOSY_MAILGUN_PUBKEY")
	mailgunPrivKey := os.Getenv("NOSY_MAILGUN_PRIVKEY")
	mailgunDomain := os.Getenv("NOSY_MAILGUN_DOMAIN")

	if mapsKey == "" || mailgunPubKey == "" || mailgunPrivKey == "" || mailgunDomain == "" {
		panic("missing required environment variable")
	}

	fromAddr := "admin@" + mailgunDomain

	toAddr := os.Getenv("NOSY_RECIPIENT")
	home := os.Getenv("NOSY_HOME")
	requestURL := os.Getenv("NOSY_311_URL")

	if toAddr == "" || home == "" || requestURL == "" {
		panic("missing required environment variable")
	}

	councilDistrict := os.Getenv("NOSY_COUNCIL_DISTRICT")

	checkIntervalStr := os.Getenv("NOSY_CHECK_INTERVAL")

	if checkIntervalStr == "" {
		checkIntervalStr = "10m"
	}

	checkInt, err := time.ParseDuration(checkIntervalStr)

	if err != nil {
		panic(err)
	}

	mg := mailgun.NewMailgun(mailgunDomain, mailgunPrivKey, mailgunPubKey)
	mapsClient, err := maps.NewClient(maps.WithAPIKey(mapsKey))

	if err != nil {
		panic(err)
	}

	return &env{
		requestURL:      requestURL,
		fromAddr:        fromAddr,
		toAddr:          toAddr,
		dist:            &googleMapsDistancer{mapsClient: mapsClient},
		mail:            &mailgunEmailer{c: mg},
		home:            home,
		councilDistrict: councilDistrict,
		checkInterval:   checkInt,
		threshold:       meters(600),
	}
}

func checkLatestRequests(e *env, startDate time.Time) error {
	latest, err := latestRequests(e.requestURL, startDate)

	if err != nil {
		return err
	}

	for _, req := range latest {
		if len(req.Geolocation.Coordinates) < 2 {
			continue
		}

		// if we've defined a council district,
		// don't even check distance if it's outside
		// the district
		if e.councilDistrict != "" {
			if req.CouncilDistrict != e.councilDistrict {
				log.Printf("skipping id %s for address %s, outside defined district", req.ServiceRequestID, req.Address)
				continue
			}
		}

		latLonPair := fmt.Sprintf("%.14f,%.14f", req.Geolocation.Coordinates[1], req.Geolocation.Coordinates[0])
		metersDist, err := e.dist.Distance(e.home, latLonPair)

		if err != nil {
			return err
		}

		log.Printf("distance for %s (id %s): %v", req.Address, req.ServiceRequestID, metersDist)

		if metersDist < e.threshold {
			e.mail.SendAlertEmail(e.fromAddr, e.toAddr, req)
		}
	}

	return nil
}
