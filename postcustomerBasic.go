package main

import (
	"github.com/gin-gonic/gin"
	//github.com/go-ozzo/ozzo-validation
	//	github.com/go-ozzo/ozzo-validation/is
)

func PostCustomerBasic(c *gin.Context) {
	/*
		getCustUniqueId := c.Query("cust_unique_id")
		lengthCustUniqueId := len(getCustUniqueId)
		if lengthCustUniqueId == 0 {
			return
		} else { //lengthCustUniqueId > 0
			flag := CheckValidCustUniqueId(getCustUniqueId)
			if flag == false {
				return
			} else { // cust_unique_id is contained in database
				fmt.Println("unique id is true ")
			} // end validid is valid

		} // end of lengthid >0
	*/ /*
		id := c.Query("id")

		first_name := c.PostForm("R1") //require
		middle_name := c.PostForm("R2")
		last_name := c.PostForm("R3") //require
		suffix := c.PostForm("R4")
		birth_date := c.PostForm("R5")
		adress1 := c.PostForm("R6")
		adress2 := c.PostForm("R7")
		country := c.PostForm("R8")
		state := c.PostForm("R9")
		city := c.PostForm("R10")
		zipcode := c.PostForm("R11")
		home_number := c.PostForm("R12")
		mobile_number := c.PostForm("R13") //require
		email := c.PostForm("R14")         //require
		program_id := c.PostForm("R15")
		cust_unique_id := c.PostForm("R16") //require

		// need to require postForm input
	*/
} // end of PstCustomerBasic function
