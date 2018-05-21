// this file is doing for CRUD
package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	//github.com/go-ozzo/ozzo-validation
	//	github.com/go-ozzo/ozzo-validation/is
)

//get infos from database base of id
func GetName(c *gin.Context) {
	id := c.Param("id")
	var Tblcustomerbasicinfo struct {
		ID        int    `db:"id" json:"id"` //json changes format print
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}

	err := db.Select("id", "first_name", "last_name").
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": id}).
		One(&Tblcustomerbasicinfo)

	fmt.Println(err)
	fmt.Println(Tblcustomerbasicinfo)
	if err != nil {
		c.JSON(404, gin.H{"error": "customer is not found"})
	}

	c.JSON(200, Tblcustomerbasicinfo) //display output in server
}

//delete infos from database of id
func DeleteName(c *gin.Context) {
	id := c.Param("id")
	tblcustomerbasicinfo := Tblcustomerbasicinfo{}

	err := db.Select("id", "first_name", "last_name").
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": id}).
		One(&tblcustomerbasicinfo)
	if err != nil {
		c.JSON(404, gin.H{"error": "customer is not found"})
	}

	db.Model(&tblcustomerbasicinfo).Delete()
	c.JSON(200, gin.H{"success": true})
}

//create a valid customer
//(if the phone and email are not the same at database,then it is valid)
func PostName(c *gin.Context) {
	firstName := c.PostForm("FirstName")
	middleName := c.PostForm("M")
	lastName := c.PostForm("LastName")
	suffix := c.PostForm("Suffix")
	birthDate := c.PostForm("BirthDate")
	address1 := c.PostForm("Address1")
	address2 := c.PostForm("Address2")

	city := c.PostForm("City")
	state := c.PostForm("State")
	country := c.PostForm("Country")
	zipCode := c.PostForm("ZipCode")
	homeNumber := c.PostForm("HomeNumber")
	mobileNumber := c.PostForm("PhoneNumber")
	email := c.PostForm("Email")
	gender := c.PostForm("Gender")
	programID := c.PostForm("ProgramID")
	programCustomerID := c.PostForm("ProgramCustomerID")
	programCustomerMobile := c.PostForm("ProgramCustomerMobile")
	custUniqueID := c.PostForm("CustUniqueID")
	createdBy := c.PostForm("CreatedBy")
	createdDate := c.PostForm("CreatedDate")
	//////////////////////////////////////////////////////////////////////////////
	//put valid check in create customer struct-------------not work properly
	/*

		validCheck := ValidationCustomer{
			valemail:        email,
			valmobileNumber: mobileNumber,
		}
		err := validCheck.validate()
		fmt.Println(err)
		if err != nil {
			c.JSON(404, gin.H{"error": "invalid email or phone"})
			//return
		}
	*/
	///////////////////////////////////////////////////////////////////////////////////
	//here need compare the phone and email with database, if one of them equal, then fail to create a customer
	/*
		var checkValidEmail struct {
			email int `db:"email" json:"Email"` //json changes format print
			//phoneNumber string `db:"mobile_number" json:"mobile_number`
		}

		err := db.Select("email").
			From("tblcustomerbasicinfo").
			Where(dbx.HashExp{"email": email}).
			One(&checkValidEmail)
		if err == nil {
			c.JSON(201, gin.H{"not valid": "you account is existed "}) //display output in server
			return
		}
	*/
	/////////////////////////////////////////////////
	//convert string to int
	fmt.Println(city)
	cityId, err := strconv.Atoi(city)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "invalid city"})
		return
	}

	stateId, err := strconv.Atoi(state)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid state"})
		return
	}

	countryId, err := strconv.Atoi(country)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid country"})
		return
	}

	eprogramID, err := strconv.Atoi(programID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid programID"})
		return
	}

	eprogramCustomerID, err := strconv.Atoi(programCustomerID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid programCustomerID"})
		return
	}

	//assign value to struct
	customerInfo := Tblcustomerbasicinfo{
		FirstName:             firstName,
		MiddleName:            middleName,
		LastName:              lastName,
		Suffix:                suffix,
		BirthDate:             birthDate,
		Address1:              address1,
		Address2:              address2,
		city:                  cityId,
		State:                 stateId,
		Country:               countryId,
		ZipCode:               zipCode,
		HomeNumber:            homeNumber,
		MobileNumber:          mobileNumber,
		Email:                 email,
		Gender:                gender,
		ProgramID:             eprogramID,
		ProgramCustomerID:     eprogramCustomerID,
		ProgramCustomerMobile: programCustomerMobile,
		CustUniqueID:          custUniqueID,
		CreatedBy:             createdBy,
		CreatedDate:           createdDate,
	}

	//fmt.Println()
	fmt.Println(db.Model(&customerInfo).Insert()) // insert struct to database table

	c.JSON(200, gin.H{" create success": firstName + lastName})
}

//////////////////////////////////////////////////////////////////////////////////
/*
func (v ValidationCustomer) validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.valemail, validation.Required, is.Email),
		validation.Field(&v.valmobileNumber, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{7-10}$"))),
	)

}
*/
///////////////////////////////////////////////////////////////////////////////////////////////
