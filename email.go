package main

import (
	"encoding/json"

	"log"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type emailer interface {
	SendErrorEmail(fromAddress string, toAddress string, err error) error
	SendAlertEmail(fromAddress string, toAddress string, serviceReq serviceRequest) error
}

type mailgunEmailer struct {
	c mailgun.Mailgun
}

func (m *mailgunEmailer) SendErrorEmail(fromAddress string, toAddress string, err error) error {
	log.Printf("sending error email: %v", err.Error())
	msg := mailgun.NewMessage(fromAddress, "NosyNeighbor Error", "Error: "+err.Error(), toAddress)

	_, _, sendErr := m.c.Send(msg)

	if sendErr != nil {
		log.Printf("Error sending error email: %v Original error: %v", sendErr, err)
	}

	return sendErr
}

func (m *mailgunEmailer) SendAlertEmail(fromAddress string, toAddress string, serviceReq serviceRequest) error {
	body, err := json.MarshalIndent(serviceReq, "", "  ")

	if err != nil {
		return m.SendErrorEmail(fromAddress, toAddress, err)
	}

	msg := mailgun.NewMessage(fromAddress, "311 Alert", string(body), toAddress)

	_, _, err = m.c.Send(msg)

	if err != nil {
		return m.SendErrorEmail(fromAddress, toAddress, err)
	}

	return nil
}
