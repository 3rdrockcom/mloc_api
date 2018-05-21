package main

//struct from database
//tblcustomerbasicinfo sturct
//create a valid customer
type Tblcustomerbasicinfo struct {
	ID                    int
	FirstName             string
	MiddleName            string
	LastName              string
	Suffix                string
	BirthDate             string
	Address1              string
	Address2              string
	city                  int
	State                 int
	Country               int
	ZipCode               string `db:"zipcode"`
	HomeNumber            string
	MobileNumber          string
	Email                 string
	Gender                string
	ProgramID             int
	ProgramCustomerID     int
	ProgramCustomerMobile string
	CustUniqueID          string
	CreatedBy             string
	CreatedDate           string
}

//////////////////////////////////////////////
/*
//for create a valid new customer
type ValidationCustomer struct {
	valemail        string
	valmobileNumber string
}
*/
//////////////////////////////////////////////
