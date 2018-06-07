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
	if len(countries) == 0 {
		err = ErrCountryNotFound
	}
	if err != nil {
		return
	}

	return
}

// GetState gets information about a state
func (ls *LookupService) GetState(countryID int, stateID int) (state models.State, err error) {
	err = DB.Select().
		From(state.TableName()).
		Where(dbx.HashExp{
			"state_id":   stateID,
			"country_id": countryID,
		}).
		One(&state)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrStateNotFound
			return
		}
		return
	}

	return
}

// GetStates gets information about every state for a country
func (ls *LookupService) GetStates(countryID int) (states models.States, err error) {
	err = DB.Select().
		From(models.State{}.TableName()).
		Where(dbx.HashExp{"country_id": countryID}).
		All(&states)
	if len(states) == 0 {
		err = ErrStateNotFound
	}
	if err != nil {
		return
	}

	return
}

// GetCity gets information about a city
func (ls *LookupService) GetCity(stateCode string, cityID int) (city models.City, err error) {
	err = DB.Select().
		From(city.TableName()).
		Where(dbx.HashExp{"state_code": stateCode}).Where(dbx.HashExp{
		"state_code": stateCode,
		"city_id":    cityID,
	}).
		One(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCityNotFound
			return
		}
		return
	}

	return
}

// GetCities gets information about every city in a state
func (ls *LookupService) GetCities(stateCode string) (cities models.Cities, err error) {
	err = DB.Select().
		From(models.City{}.TableName()).
		Where(dbx.HashExp{"state_code": stateCode}).
		All(&cities)
	if len(cities) == 0 {
		err = ErrCityNotFound
	}
	if err != nil {
		return
	}

	return
}

// GetIncomeSource gets information about an income source
func (ls *LookupService) GetIncomeSource(id int) (incomeSource models.IncomeSource, err error) {
	err = DB.Select().
		From(incomeSource.TableName()).
		Where(dbx.HashExp{"id": id}).
		One(&incomeSource)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrIncomeSourceNotFound
			return
		}
		return
	}

	return
}

// GetIncomeSources gets information about every income source
func (ls *LookupService) GetIncomeSources() (incomeSources models.IncomeSources, err error) {
	err = DB.Select().
		From(models.IncomeSource{}.TableName()).
		All(&incomeSources)
	if len(incomeSources) == 0 {
		err = ErrIncomeSourceNotFound
	}
	if err != nil {
		return
	}

	return
}

// GetPayFrequency gets information about a pay frequency
func (ls *LookupService) GetPayFrequency(id int) (payFrequency models.PayFrequency, err error) {
	err = DB.Select().
		From(payFrequency.TableName()).
		Where(dbx.HashExp{"id": id}).
		One(&payFrequency)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrPayFrequencyNotFound
			return
		}
		return
	}

	return
}

// GetPayFrequencies gets information about every pay frequency
func (ls *LookupService) GetPayFrequencies() (payFrequencies models.PayFrequencies, err error) {
	err = DB.Select().
		From(models.PayFrequency{}.TableName()).
		All(&payFrequencies)
	if len(payFrequencies) == 0 {
		err = ErrPayFrequencyNotFound
	}
	if err != nil {
		return
	}

	return
}
