package customer

import (
	"github.com/epointpayment/mloc_api_go/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Settings manages customer setting information
type Settings struct {
	cs *CustomerService
}

// Get gets a global setting entry
func (s *Settings) Get(id int) (setting models.SystemSetting, err error) {
	err = DB.Select().
		From(setting.TableName()).
		Where(dbx.HashExp{"id": id}).
		One(&setting)
	if err != nil {
		err = ErrProblemOccured
		return
	}

	return
}
