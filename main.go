// main function to run crud
// run lookup api
//run registration api
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"

	_ "github.com/go-sql-driver/mysql"
)

var db *dbx.DB

//main function to call RCUD,lookup API and registration API
func main() {

	router := gin.Default()
	//connect to localhost database
	var err error
	db, err = dbx.Open("mysql", "root:@tcp(127.0.0.1:3306)/mloc_live") // connect database
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.LogFunc = log.Printf

	/**************router********************************/
	//crud
	router.GET("/user/get/:id", GetName)  //countryId has a value
	router.POST("/user/create", PostName) //countryId has a value
	//router.PUT("/user/update/:name", putName)			//countryId has a value
	router.DELETE("/user/delete/:id", DeleteName) //function is in controllers.go file
	//look up API
	router.GET("/api/v1/lookup/get_country/", GetCountryId)          // function is in lookupContry.go file
	router.GET("/api/v1/lookup/get_state/", GetStateId)              //function is in lookupCity.go file
	router.GET("/api/v1/lookup/get_city/", GetCity)                  //function is in lookupState.go file
	router.GET("/api/v1/lookup/get_income_source/", GetIncomeSource) // function is in lookupIncomeSource.go file
	router.GET("/api/v1/lookup/get_pay_frequency/", GetPayFrequency) // function is in lookupPayFrequency.go file

	//registration API
	//	router.Get("/api/v1/login/get_customer_key/", GetCustomerKey)
	router.GET("/api/v1/customer/get_customer/", GetCustomer)                      // funciton is in getCustomer.go file
	router.GET("/customer/get_customer_loan/", GetCustomerLoan)                    //function is in getCustomerLoan.go file
	router.GET("/api/v1/customer/get_transaction_history/", GetTransactionHistory) //fuction is in getTransactionHistory.go file

	//post
	router.POST("/api/v1/customer/customer_basic/", PostCustomerBasic)
	router.POST("/api/v1/customer/customer_additional/", PostCustomerAdditional)
	router.POST("/api/v1/customer/accept_terms_and_condition/", PostCreditLineApplication)
	router.Run(":8080")
}
