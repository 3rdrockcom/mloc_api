package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

//this function is to get customerloan
// if the query is empty, display all loan of customer
// if the query is not empty, then check it is valid or not
// if the query is valid, display info
// if the query is not valid, display info of error

func GetCustomerLoan(c *gin.Context) {
	getQueryCustomerLoan := c.Query("R1") // get query
	if len(getQueryCustomerLoan) == 0 {   // empty query, display all loan of customer
		return
	} else {
		flag := CheckValidCustUniqueId(getQueryCustomerLoan)
		//invalid query of cust_unique_id
		if flag == false {
			c.JSON(404, gin.H{"status": false, "message": "No customer were found.", "response_code": 404})
		} else {
			// valid query in database
			var TblGetCustomerLoan struct {
				ID                   *string `db:"loanId" json"id"`
				FKCustomerId         *string `db:"fk_customer_id" json"fk_customer_id"`
				TotalPrincipalAmount *string `db:"total_principal_amount" json:"total_principal_amount"`
				TotalFeeAmount       *string `db"total_fee_amount" json:"total_fee_amount"`
				TotalAmount          *string `db:"total_amount" json:"total_amount"`
			} //end of TblGetCustoemrLoan struct
			err := db.Select("tblcustomerloantotal.id As loanId", "tblcustomerloantotal.fk_customer_id", "tblcustomerloantotal.total_principal_amount", "tblcustomerloantotal.total_fee_amount", "tblcustomerloantotal.total_amount").
				From("tblcustomerbasicinfo").
				LeftJoin("tblcustomerotherinfo", dbx.NewExp("tblcustomerotherinfo.fk_customer_id = tblcustomerbasicinfo.id")).
				//LeftJoin("tblcustomerloantotal", dbx.NewExp("tblcustomerloantotal.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcountry", dbx.NewExp("tblcountry.country_id = tblcustomerbasicinfo.country")).
				LeftJoin("tblstate", dbx.NewExp("tblstate.state_id =tblcustomerbasicinfo.state ")).
				LeftJoin("tblcity", dbx.NewExp("tblcity.city_id =tblcustomerbasicinfo.city ")).
				LeftJoin("tblincomesource", dbx.NewExp("tblincomesource.id =tblcustomerotherinfo.income_source ")).
				LeftJoin("tblpayfrequency", dbx.NewExp("tblpayfrequency.id =tblcustomerotherinfo.pay_frequency ")).
				LeftJoin("tblcustomeragreement", dbx.NewExp("tblcustomeragreement.fk_customer_id =tblcustomerbasicinfo.id ")).
				LeftJoin("tblapikey", dbx.NewExp("tblapikey.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcustomercreditline", dbx.NewExp("tblcustomercreditline.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcustomerloantotal", dbx.NewExp("tblcustomerloantotal.fk_customer_id = tblcustomerbasicinfo.id")).
				Where(dbx.HashExp{"cust_unique_id": getQueryCustomerLoan}).
				One(&TblGetCustomerLoan) //fetch a single row in struct n

			fmt.Println(err)
			if err != nil {
				c.JSON(404, gin.H{"status": false, "message": "No customer were found.", "response_code": 404})

			}
			c.JSON(200, TblGetCustomerLoan)

		}
	} // end of else leng of query >0

} //end of function
