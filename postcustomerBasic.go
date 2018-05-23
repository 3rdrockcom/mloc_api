package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type updateCustomerBasic struct {
	Id                 int
	UpdateFirstName    *string `db:"first_name"`
	UpdateMiddleName   *string `db:"middle_name"`
	UpdateLastName     *string `db:"last_name"`
	UpdateMobileNumber *string `db:"mobile_number"`
	UpdateEmail        *string `db:"email"`
}

func (c updateCustomerBasic) TableName() string {
	return "tblcustomerbasicinfo"
}

func (a updateCustomerBasic) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UpdateFirstName, validation.Required),
		validation.Field(&a.UpdateLastName, validation.Required),
		validation.Field(&a.UpdateMobileNumber, validation.Required),
		validation.Field(&a.UpdateEmail, validation.Required, is.Email),
	)
}

func PostCustomerBasic(c *gin.Context) {
	tempId := 15
	updateinfo := &updateCustomerBasic{}
	err := db.Select().
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"id": tempId}).
		One(updateinfo)
	if err != nil {
		c.JSON(404, gin.H{"error": "404 Not Found"})

	}
	fmt.Println(updateinfo)
	formKeys := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15", "R16"}
	//c.GetPostForm(FormKeys)

	for index := range formKeys {
		formKey := formKeys[index]
		//fmt.Println(c.PostForm(formKey))
		value, isNotNull := c.GetPostForm(formKey)
		if isNotNull {

			fmt.Println(formKey + " " + value)
			switch formKey {
			case "R1":
				updateinfo.UpdateFirstName = &value
			case "R2":
				updateinfo.UpdateMiddleName = &value

			case "R3":
				updateinfo.UpdateLastName = &value
			case "R13":
				updateinfo.UpdateMobileNumber = &value
			case "R14":
				updateinfo.UpdateEmail = &value

			}
		}
	}

	//check require postform is valid or not
	err = updateinfo.Validate()

	if err != nil {
		c.JSON(400, gin.H{"status": false, "message": "Provide complete customer information to create.", "response_code": 400}) // invalid postformretur
		return
	}

	err = db.Model(updateinfo).Update()
	if err != nil {
		c.JSON(500, gin.H{"status": false, "message": "Provide complete customer information to create.", "response_code": 400}) // invalid postformretur
		return
	} // end of valid postform
	c.JSON(200, gin.H{"status": true, "message": "customer information has been updated successfully.", "response_code": 200})

} // end of PstCustomerBasic function
