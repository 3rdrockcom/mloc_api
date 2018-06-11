package notifications

import (
	"github.com/epointpayment/mloc_api_go/app/config"

	Mail "github.com/epointpayment/mloc_api_go/app/services/notifications/mail"
	SMS "github.com/epointpayment/mloc_api_go/app/services/notifications/sms"
)

// NotificationsService is a service that manages notifications
type NotificationsService struct{}

// New creates an instance of the service
func New() *NotificationsService {
	return &NotificationsService{}
}

func (ns *NotificationsService) Send(payload interface{}) (err error) {
	switch payload.(type) {
	case *Mail.Mail:
		// USe SMTP
		if config.Get().Mail.Driver == "smtp" {
			n := payload.(*Mail.Mail)
			err = n.Send()
			if err != nil {
				return
			}
			return
		}
	case *SMS.SMS:
		// Use InfoBIP
		if config.Get().SMS.Driver == "infobip" {
			n := payload.(*SMS.SMS)
			_, err = n.Send()
			if err != nil {
				return
			}
			return
		}
	default:
		err = ErrInvalidPayloadType
		return
	}

	return
}
