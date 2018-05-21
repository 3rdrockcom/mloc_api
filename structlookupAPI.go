package main

//table of tblcountry from database
//declair struct for lookup api

type Tblcountry struct {
	CountryId        int
	Name             string
	IsoCode2         string
	IsoCode3         string
	AddressFormat    string
	PostCodeRequired int
	Status           int
	MobilePrefix     string
}

type Tblstate struct {
	StateId   int
	State     string
	StateCode string
	CountryId int
}

type Tblcity struct {
	CityId    int
	City      string
	StateCode string
}

type TblIncomeSource struct {
	Id          int
	Description string
}

type TblPayFrequency struct {
	Id          int
	Description string
}
