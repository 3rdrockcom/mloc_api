package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

//below function is Get get_customer
func GetCustomer(c *gin.Context) {
	queryCustUniqueId := c.Query("cust_unique_id") //get the query
	//check the query is valid or not
	if len(queryCustUniqueId) == 0 {
		return
	} else { //else len of queryCustUniqueId is > 0

		flag := CheckValidCustUniqueId(queryCustUniqueId)
		//	check the query if it is not in database
		if flag == false {
			c.JSON(404, gin.H{"status": false, "message": "No customer were found.", "response_code": 404})
		} else {
			//query is valid
			//fetch a single row in struct
			var TblGetCustomer struct {
				CustUniqueId     string  `db:"cust_unique_id" json:"cust_unique_id"`
				Id               *string `db:"id" json:"customer_id"`
				FirstName        *string `db:"first_name" json:"first_name"`
				MiddleName       *string `db:"middle_name" json:"middle_name"`
				LastName         *string `db:"last_name" json:"last_name"`
				Suffix           *string `db:"suffix" json:"suffix"`
				BirthDate        *string `db:"birth_date" json:"birth_date"`
				Address1         *string `db:"address1" json:"address1"`
				Address2         *string `db:"address2" json:"address2"`
				CountryId        *string `db:"country_id" json:"country_id"`
				CountryDesc      *string `db:"name" json:"country_desc"`
				StateId          *string `db:"state_id" json:"state_id"`
				StateDesc        *string `db:"state" json:"state_desc"`
				CityId           *string `db:"city_id" json:"city_id"`
				CityDesc         *string `db:"city" json:"city_desc'`
				ZipCode          *string `db:"zipcode" json:"zipcode"`
				HomeNumber       *string `db:"home_number" json:"home_number"`
				MobileNumber     *string `db:"mobile_number" json: "mobile_number"`
				Email            *string `db:"email" json:"email"`
				CompanyName      *string `db:"company_name" json:"company_name"`
				PhoneNumber      *string `db:"phone_number" json:"phone_number"`
				NetPayPerCheck   *string `db:"net_pay_percheck" json:"net_pay_percheck"`
				IncomeSourceId   *string `db:"income_source" json:"income_source_id"`
				MlocAccess       *string `db:"mloc_access" json:"mloc_access"`
				Registration     *string `db:"registration" json:"registration"`
				TermAndCondition *string `db:"term_and_condition" json:"term_and_condition"`
				IncomeSourceDesc *string `db:"description" json:"income_source_desc"`
				PayFrequencyId   *string `db:"payfreId" json:"pay_frequency_id"`
				PayFrequencyDesc *string `db:"payfreDesc"" json:"pay_frequency_desc"`
				NextPayDate      *string `db:"next_paydate" json:"next_paydate"`
				Key              *string `db:"key" json:"key"`
				CreditLimit      *string `db:"credit_limit" json:"credit_limit"`
				AvailableCredit  *string `db:"available_credit" json:"available_credit"`
				IsSuspended      string  `db:"is_suspended"  json:"is_suspended"`
			} // end varialbe

			//building select query
			err := db.Select("tblcustomerbasicinfo.cust_unique_id", "tblcustomerbasicinfo.id", "first_name", "middle_name", "last_name", "tblcustomerbasicinfo.suffix",
				"birth_date", "address1", "address2", "tblcountry.country_id", "tblcountry.name",
				"tblstate.state_id", "tblstate.state", "tblcity.city_id", "tblcity.city", "zipcode", "home_number", "mobile_number", "email",
				"company_name", "phone_number", "tblcustomerotherinfo.net_pay_percheck", "tblcustomerotherinfo.income_source", "tblcustomeragreement.mloc_access", "tblcustomeragreement.registration",
				"tblcustomeragreement.term_and_condition", "tblincomesource.description", "tblpayfrequency.id AS payfreId", "tblpayfrequency.description AS payfreDesc", "tblcustomerotherinfo.next_paydate",
				"key", "credit_limit", "available_credit", "is_suspended").
				From("tblcustomerbasicinfo").
				LeftJoin("tblcustomerotherinfo", dbx.NewExp("tblcustomerotherinfo.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcountry", dbx.NewExp("tblcountry.country_id = tblcustomerbasicinfo.country")).
				LeftJoin("tblstate", dbx.NewExp("tblstate.state_id =tblcustomerbasicinfo.state ")).
				LeftJoin("tblcity", dbx.NewExp("tblcity.city_id =tblcustomerbasicinfo.city ")).
				LeftJoin("tblincomesource", dbx.NewExp("tblincomesource.id =tblcustomerotherinfo.income_source ")).
				LeftJoin("tblpayfrequency", dbx.NewExp("tblpayfrequency.id =tblcustomerotherinfo.pay_frequency ")).
				LeftJoin("tblcustomeragreement", dbx.NewExp("tblcustomeragreement.fk_customer_id =tblcustomerbasicinfo.id ")).
				LeftJoin("tblapikey", dbx.NewExp("tblapikey.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcustomercreditline", dbx.NewExp("tblcustomercreditline.fk_customer_id = tblcustomerbasicinfo.id")).
				LeftJoin("tblcustomerloantotal", dbx.NewExp("tblcustomerloantotal.fk_customer_id = tblcustomerbasicinfo.id")).
				Where(dbx.HashExp{"cust_unique_id": queryCustUniqueId}).
				One(&TblGetCustomer) //fetch a single row in struct
			fmt.Println(err)
			if err != nil {
				c.JSON(404, gin.H{"status": false, "message": "No customer were found.", "response_code": 404})

			}
			//set IsSuspended to interge from string
			if TblGetCustomer.IsSuspended == "NO" {
				TblGetCustomer.IsSuspended = "0"
				c.JSON(200, TblGetCustomer)
			} else {
				TblGetCustomer.IsSuspended = "1"
				c.JSON(200, TblGetCustomer)
			}
		}
	} // end else if that has a query

} // end function
