package notifications

import (
	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/labstack/gommon/log"

	"github.com/nuveo/infobip"
)

// cfg caches the config
var cfg config.SMS

// client caches the client
var client *infobip.Client

// SMS is a service that manages the sms provider
type SMS struct {
	From string
	To   string
	Body string
}

// New creates an instance of the sms service
func New() *SMS {
	// Initialize config and client
	if client == nil {
		cfg = config.Get().SMS
		client = infobip.ClientWithBasicAuth(cfg.Username, cfg.Password)
	}

	return &SMS{
		From: cfg.FromName,
	}
}

// Send sends the request to the sms provider for processing
func (sms *SMS) Send() (response infobip.Response, err error) {
	// From
	from := cfg.FromName
	if sms.From != "" {
		from = sms.From
	}

	// To
	to := sms.format(sms.To)
	if config.IsDev() && cfg.ToMobile != "" {
		to = sms.format(cfg.ToMobile)
	}

	// Prepare payload
	m := infobip.Message{
		From: from,
		To:   to,
		Text: sms.Body,
	}
	err = m.Validate()
	if err != nil {
		return
	}

	// Send payload to sms provider
	response, err = client.SingleMessage(m)
	if err != nil {
		return
	}

	// Determine if sms provider was successful in processing request
	for _, message := range response.Messages {
		if message.Status.GroupID != 1 {
			// Log response if status group id had issues
			log.Warn(response)
			return
		}
	}

	return
}

// format formats the phone number
func (sms *SMS) format(phoneNumber string) (phoneNumberFormatted string) {
	phoneNumberFormatted = phoneNumber

	if len(phoneNumber) == 10 {
		phoneNumberFormatted = "1" + phoneNumber
	}

	return
}
