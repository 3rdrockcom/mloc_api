package main

import (
	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func GetTransactionHistory(c *gin.Context) {

	customerId := 16
	getCustUniqueId := c.Query("R1") // get first query
	firstQueryLength := len(getCustUniqueId)

	if firstQueryLength == 0 { // if the first query is empty,then return
		return
	} else if firstQueryLength > 0 {
		flag := CheckValidCustUniqueId(getCustUniqueId)
		//fmt.Println(flag)
		if flag == false {
			c.JSON(404, gin.H{"status": false, "message": "No customer were found.", "response_code": 404})
		} else {

			// there is valid cust unique id
			getLoanOrPayment := c.Query("R2") //get second query
			if len(getLoanOrPayment) > 0 {    //second query contains string
				if getLoanOrPayment == "LOAN" || getLoanOrPayment == "PAYMENT" { //second query is LOAN string
					q := db.Select().
						From("view_transaction_history").
						//Where(dbx.Like("t_type", getLoanOrPayment)).
						Where(dbx.HashExp{"t_type": getLoanOrPayment, "fk_customer_id": customerId}). // the 16 is instann , i need to get id from tblcustomerbasicinfo
						OrderBy("fk_customer_id")

					var Tblcustomerloanhistory []struct {
						FkCustomerId string  `db:"fk_customer_id" json:"fk_customer_id"`
						Amount       *string `db:"amount" json:"amount"`
						Ttype        *string `db:"t_type" json:"t_type"`
						TDate        *string `db:"t_date" json:"t_date"`
					}
					q.All(&Tblcustomerloanhistory)
					c.JSON(200, Tblcustomerloanhistory)

				} else { // second query is not LOAN and PAYMENT string
					return
				}

			} else if len(getLoanOrPayment) == 0 { //second query is empty,display all loan and payment

				//q := db.NewQuery("SELECT * FROM view_transaction_history") // call aTblcustomerloanhistory struct
				q := db.Select().
					From("view_transaction_history").
					Where(dbx.HashExp{"fk_customer_id": customerId}).
					OrderBy("fk_customer_id")

				var Tblcustomerloanhistory []struct {
					FkCustomerId string  `db:"fk_customer_id" json:"fk_customer_id"`
					Amount       *string `db:"amount" json:"amount"`
					Ttype        string  `db:"t_type" json:"t_type"`
					TDate        *string `db:"t_date" json:"t_date"`
				}
				q.All(&Tblcustomerloanhistory)
				c.JSON(200, Tblcustomerloanhistory)

			} else { //second query is nothing
				return
			}

		} // end of valid customer unique id
	} // end of first query contains value
} // end of function
