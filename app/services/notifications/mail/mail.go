package notifications

import (
	"github.com/epointpayment/mloc_api_go/app/config"

	gomail "github.com/go-mail/mail"
)

// cfg caches the config
var cfg config.Mail

// client caches the client
var client *gomail.Dialer

// Mail is a service that manages the mail provider
type Mail struct {
	From        Address
	To          []Address
	CC          []Address
	BCC         []Address
	Subject     string
	Body        string
	BodyHTML    string
	Attachments []File
	Embeds      []File
}

// Address
type Address struct {
	Address string
	Name    string
}

// File
type File struct {
	File     string
	RenameTo string
}

// New creates an instance of the mail service
func New() *Mail {
	// Initialize config and client
	if client == nil {
		cfg = config.Get().Mail
		client = gomail.NewDialer(cfg.Host, int(cfg.Port), cfg.Username, cfg.Password)
	}

	return &Mail{}
}

// Send
func (m *Mail) Send() (err error) {
	addresses := []string{}

	// Initialize message
	msg := gomail.NewMessage()

	// From
	from := Address{Address: cfg.FromAddress, Name: cfg.FromName}
	if m.From.Address != "" {
		from.Address = m.From.Address
		from.Name = m.From.Name
	}
	msg.SetHeader("From", msg.FormatAddress(from.Address, from.Name))

	// To
	addresses = m.formatAddresses(msg, m.To)
	if config.IsDev() && cfg.ToAddress != "" {
		addresses = m.formatAddresses(msg, []Address{Address{Address: cfg.ToAddress}})
	}
	msg.SetHeader("To", addresses...)

	// CC
	if len(m.CC) > 0 {
		addresses = m.formatAddresses(msg, m.CC)
		msg.SetHeader("Cc", addresses...)
	}

	// BCC
	if len(m.BCC) > 0 {
		addresses = m.formatAddresses(msg, m.BCC)
		msg.SetHeader("Bcc", addresses...)
	}

	// Subject
	msg.SetHeader("Subject", m.Subject)

	// Body - PlainText
	msg.SetBody("text/plain", m.Body)

	// Body - HTML
	if m.BodyHTML == "" {
		m.BodyHTML = "<html><body><div>" + m.Body + "</div></body></html>"
	}
	msg.AddAlternative("text/html", m.BodyHTML)

	// Attachments
	for _, attachment := range m.Attachments {
		if attachment.RenameTo == "" {
			msg.Attach(attachment.File)
			continue
		}

		msg.Attach(attachment.File, gomail.Rename(attachment.RenameTo))
	}

	// Embeds
	for _, embed := range m.Embeds {
		if embed.RenameTo == "" {
			msg.Embed(embed.File)
			continue
		}

		msg.Embed(embed.File, gomail.Rename(embed.RenameTo))
	}

	// Send payload to sms provider
	err = client.DialAndSend(msg)
	if err != nil {
		return
	}

	return
}

// formatAddresses formats the email addresses
func (m *Mail) formatAddresses(msg *gomail.Message, addresses []Address) (list []string) {
	for _, address := range addresses {
		list = append(list, msg.FormatAddress(address.Address, address.Name))
	}

	return
}
