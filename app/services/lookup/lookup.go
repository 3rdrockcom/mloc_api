package lookup

import (
	"database/sql"

	"github.com/epointpayment/mloc_api_go/app/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DB is the database handler
var DB *dbx.DB

// LookupService is a service that manages the API access
type LookupService struct{}

// New creates an instance of the service
func New() *LookupService {
	return &LookupService{}
}

// GetCountry gets information about a country
func (ls *LookupService) GetCountry(countryID int) (country models.Country, err error) {
	err = DB.Select().
		From(country.TableName()).
		Where(dbx.HashExp{"country_id": countryID}).
		One(&country)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCountryNotFound
			return
		}
		return
	}

	return
}

// GetCountries gets information about every country
func (ls *LookupService) GetCountries() (countries models.Countries, err error) {
	err = DB.Select().
		From(models.Country{}.TableName()).
		All(&countries)
	if err != nil {
		return
	}

	return
}
