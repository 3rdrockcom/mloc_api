package payments

import (
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/disbursement"
	"github.com/epointpayment/mloc_api_go/app/services/payments/driver/epoint"
)

// PaymentsService is a service that manages payment disbursements and collections
type PaymentsService struct{}

// New creates an instance of the service
func New() *PaymentsService {
	return &PaymentsService{}
}

func (s *PaymentsService) Disbursement(req disbursement.Request) (res disbursement.Response, err error) {
	switch req.Method {
	case "epoint":
		var es *epoint.Driver
		es, err = epoint.New()
		if err != nil {
			return
		}
		res, err := es.Disbursement(req)
		if err != nil {
			return res, err
		}

	default:
		err = ErrInvalidPayloadType
		return
	}

	return
}

func (s *PaymentsService) Collection(req collection.Request) (res collection.Response, err error) {
	switch req.Method {
	case "epoint":
		var es *epoint.Driver
		es, err = epoint.New()
		if err != nil {
			return
		}
		res, err := es.Collection(req)
		if err != nil {
			return res, err
		}

	default:
		err = ErrInvalidPayloadType
		return
	}

	return
}
