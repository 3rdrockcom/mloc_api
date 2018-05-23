package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type UpdateCustAdditional struct {
	Id               int
	FkCustomerId     int      `db:"fk_customer_id"`
	CompanyName      *string  `db:"company_name"`
	PhoneNumber      *string  `db:"phone_number"`
	NetPayPercheck   *float64 `db:"net_pay_percheck"`
	IncomeSource     *int     `db:"income_source"`
	Payfrequency     *int     `db:"pay_frequency"`
	NextPayDate      *string  `db:"next_paydate"`
	FollowingPayDate *string  `db:"following_paydate"`
	//CustUniqueId *string `db:"cust_unique_id`
}

func (c UpdateCustAdditional) TableName() string {
	return "tblcustomerotherinfo"
}

func PostCustomerAdditional(c *gin.Context) {
	tempFkCustId := 15
	updateCustomerOtherinfo := &UpdateCustAdditional{}
	err := db.Select().
		From("tblcustomerotherinfo").
		Where(dbx.HashExp{"Fk_customer_id": tempFkCustId}).
		One(updateCustomerOtherinfo)
	if err != nil {
		c.JSON(404, gin.H{"error": "404 Not Found"})

	}

	formKeys := []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8"}

	for index := range formKeys {
		formKey := formKeys[index]
		value, isNotNull := c.GetPostForm(formKey)
		if isNotNull {

			fmt.Println(formKey + " " + value)
			switch formKey {
			case "R1":
				updateCustomerOtherinfo.CompanyName = &value
			case "R2":
				updateCustomerOtherinfo.PhoneNumber = &value
			case "R3":
				fvalue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					fmt.Println(err)
					c.JSON(400, gin.H{"status": false, "message": "Provide complete customer information to create.", "response_code": 400})
				}
				updateCustomerOtherinfo.NetPayPercheck = &fvalue
			case "R4":
				var tempVal int
				tempVal, err = strconv.Atoi(value)
				if err != nil {
					fmt.Println(err)
					c.JSON(400, gin.H{"status": false, "message": "Provide complete customer information to create.", "response_code": 400})
					return
				}
				updateCustomerOtherinfo.IncomeSource = &tempVal

			case "R5":
				var tempVal int
				tempVal, err = strconv.Atoi(value)
				if err != nil {
					fmt.Println(err)
					c.JSON(400, gin.H{"status": false, "message": "Provide complete customer information to create.", "response_code": 400})
					fmt.Println(err)
					return
				}
				updateCustomerOtherinfo.Payfrequency = &tempVal
			case "R6":
				updateCustomerOtherinfo.NextPayDate = &value
			case "R7":
				updateCustomerOtherinfo.FollowingPayDate = &value
				//case "R8":
				//	updateCustomerOtherinfo.CustUniqueId = &value

			} //end switch

		} else {
			c.JSON(500, gin.H{"status": false, "message": "Provide complete customer information to create", "response_code": 400}) // invalid postformretur
			fmt.Println(err)
		}
	}
	err = db.Model(updateCustomerOtherinfo).Update()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"status": false, "message": "Some problems occurred,please try again.", "response_code": 400}) // invalid postformretur

		fmt.Println(err)
		return
	} // end of valid postform
	c.JSON(200, gin.H{"status": true, "message": "customer information has been updated successfully.", "response_code": 200})

}
