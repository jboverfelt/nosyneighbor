# nosyneighbor [![builds.sr.ht status](https://builds.sr.ht/~jboverfelt/nosyneighbor.svg)](https://builds.sr.ht/~jboverfelt/nosyneighbor?)

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
	* `NOSY_MAPS_KEY` - the Google Maps API key (required)
	* `NOSY_MAILGUN_PUBKEY` - the public key for your Mailgun account (required)
	* `NOSY_MAILGUN_PRIVKEY` - the private key for your Mailgun account (required)
	* `NOSY_MAILGUN_DOMAIN` - the domain from which emails will be sent (must be a domain you've configured in Mailgun, required)
	* `NOSY_RECIPIENT` - the email address to send alert emails to (required)
	* `NOSY_311_URL` - the Socrata webservice that houses the 311 requests (ex https://data.baltimorecity.gov/resource/ni4d-8w7k.json) (required)
	* `NOSY_CHECK_INTERVAL` - the time interval to check for new 311 requests. Accepts the format used by time.ParseDuration. Default is "10m" if not set
	* `NOSY_HOME` - your home address
	* `NOSY_COUNCIL_DISTRICT` - Optional. Setting this filters out any requests outside of the specified City Council District. This is used to reduce the number of requests to the Distance API
3. `$ ./nosyneighbor`
