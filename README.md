# nosyneighbor

Get email notifications when 311 requests are opened near your house

nosyneighbor is intended to be run on a server. This can just be a computer
on your home network that has access to the internet. This computer does not need
to be accessible from outside your network. 

## Prerequisites

* Go 1.7 or later (will probably work with earlier versions, but untested)
* A Google account and API key
* A Mailgun account and API key

Note: this was built and tested on Linux, but no reason it shouldn't work on any other platform Go supports

## Compiling and Running

1. All dependencies are in the vendor folder, just `go build`
2. Define the following environment variables:
    * `NOSY_MAPS_KEY` - the Google Maps API key
    * `NOSY_MAILGUN_PUBKEY` - the public key for your Mailgun account
    * `NOSY_MAILGUN_PRIVKEY` - the private key for your Mailgun account
    * `NOSY_MAILGUN_DOMAIN` - the domain from which emails will be sent (must be a domain you've configured in Mailgun)
    * `NOSY_RECIPIENT` - the email address to send alert emails to
    * `NOSY_311_URL` - the Open311 compliant "requests" webservice (ex http://311.baltimorecity.gov/open311/v2/requests.json)
3. `$ ./nosyneighbor`