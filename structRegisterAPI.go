package main

//declair TblGetCustomer struct for get customer of registration api in getcustomer.go file
type TblGetCustomer struct {
	CustUniqueId     string
	Id               string
	FirstName        string
	MiddleName       string
	LastName         string
	Suffix           *string
	BirthDate        string
	Address1         *string
	Address2         *string
	CountryId        *string
	CountryDesc      *string
	StateId          *string
	StateDesc        *string
	CityId           *string
	CityDesc         *string
	ZipCode          *string
	HomeNumber       *string
	MobileNumber     *string
	Email            *string
	CompanyName      *string
	PhoneNumber      *string
	NetPayPerCheck   *string
	IncomeSourceId   string
	MlocAccess       *string
	Registration     *string
	TermAndCondition *string
	IncomeSourceDesc *string
	PayFrequencyId   string
	PayFrequencyDesc *string
	NextPayDate      *string
	Key              string
	CreditLimit      *string
	AvailableCredit  *string
	IsSuspended      string
} // end varialbe

type TblGetCustomerLoan struct {
	ID                   *string
	FKCustomerId         *string
	TotalPrincipalAmount *string
	TotalFeeAmount       *string
	TotalAmount          *string
}
