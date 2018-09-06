package payments

import (
	"github.com/epointpayment/mloc_api_go/app/services/payments/collection"
	"github.com/epointpayment/mloc_api_go/app/services/payments/disbursement"
	"github.com/epointpayment/mloc_api_go/app/services/payments/driver/epoint"
	"github.com/epointpayment/mloc_api_go/app/services/payments/driver/stp"
	"github.com/epointpayment/mloc_api_go/app/services/payments/registration"
)

var paymentDisbursementMethod = make(map[string]PaymentDisbursementMethod)
var paymentCollectionMethod = make(map[string]PaymentCollectionMethod)

func init() {
	disbursementMethods := []PaymentDisbursementMethod{
		PaymentDisbursementMethod{ID: 1, Code: MethodEPOINT, Name: "Wallet"},
		PaymentDisbursementMethod{ID: 2, Code: MethodSTP, Name: "Bank Account"},
	}
	for _, method := range disbursementMethods {
		paymentDisbursementMethod[method.Code] = method
	}

	collectionMethods := []PaymentCollectionMethod{
		PaymentCollectionMethod{ID: 1, Code: MethodEPOINT, Name: "Wallet"},
		PaymentCollectionMethod{ID: 2, Code: MethodSTP, Name: "Bank Account"},
	}
	for _, method := range collectionMethods {
		paymentCollectionMethod[method.Code] = method
	}
}

// PaymentsService is a service that manages payment disbursements and collections
type PaymentsService struct{}

// New creates an instance of the service
func New() *PaymentsService {
	return &PaymentsService{}
}

func (s *PaymentsService) Register(req registration.Request) (res registration.Response, err error) {
	switch req.Method {
	case MethodSTP:
		var driver *stp.Driver
		driver, err = stp.New()
		if err != nil {
			return
		}
		res, err = driver.Register(req)
		if err != nil {
			return res, err
		}
	default:
		err = ErrInvalidPayloadType
		return
	}

	return
}

func (s *PaymentsService) Disbursement(req disbursement.Request) (res disbursement.Response, err error) {
	switch req.Method {
	case MethodEPOINT:
		var es *epoint.Driver
		es, err = epoint.New()
		if err != nil {
			return
		}
		res, err = es.Disbursement(req)
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
	case MethodEPOINT:
		var es *epoint.Driver
		es, err = epoint.New()
		if err != nil {
			return
		}
		res, err = es.Collection(req)
		if err != nil {
			return res, err
		}

	default:
		err = ErrInvalidPayloadType
		return
	}

	return
}

type PaymentDisbursementMethod struct {
	ID   int    `json:"-"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type PaymentCollectionMethod struct {
	ID   int    `json:"-"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (s *PaymentsService) GetDisbursementMethods() (entries []PaymentDisbursementMethod, err error) {
	for _, method := range paymentDisbursementMethod {
		entries = append(entries, method)
	}

	return
}

func (s *PaymentsService) GetCollectionMethods() (entries []PaymentCollectionMethod, err error) {
	for _, method := range paymentCollectionMethod {
		entries = append(entries, method)
	}

	return
}
