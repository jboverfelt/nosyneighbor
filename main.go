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
	requestURL    string
	fromAddr      string
	toAddr        string
	dist          distancer
	mail          emailer
	home          string
	threshold     meters
	checkDuration time.Duration
}

func main() {
	curEnv := setupEnv()

	for {
		startDate := time.Now().Add(-curEnv.checkDuration)
		err := checkLatestRequests(curEnv, startDate)

		if err != nil {
			curEnv.mail.SendErrorEmail(curEnv.fromAddr, curEnv.toAddr, err)
			return
		}

		time.Sleep(curEnv.checkDuration)
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

	mg := mailgun.NewMailgun(mailgunDomain, mailgunPrivKey, mailgunPubKey)
	mapsClient, err := maps.NewClient(maps.WithAPIKey(mapsKey))

	if err != nil {
		panic(err)
	}

	return &env{
		requestURL:    requestURL,
		fromAddr:      fromAddr,
		toAddr:        toAddr,
		dist:          &googleMapsDistancer{mapsClient: mapsClient},
		mail:          &mailgunEmailer{c: mg},
		home:          home,
		checkDuration: 2 * time.Minute,
		threshold:     meters(600),
	}
}

func checkLatestRequests(e *env, startDate time.Time) error {
	latest, err := latestRequests(e.requestURL, startDate)

	if err != nil {
		return err
	}

	for _, req := range latest {
		latLonPair := fmt.Sprintf("%f,%f", req.Lat, req.Long)
		metersDist, err := e.dist.Distance(e.home, latLonPair)

		if err != nil {
			return err
		}

		log.Printf("distance for %s: %v", req.Address, metersDist)

		if metersDist < e.threshold {
			e.mail.SendAlertEmail(e.fromAddr, e.toAddr, req)
		}
	}

	return nil
}
