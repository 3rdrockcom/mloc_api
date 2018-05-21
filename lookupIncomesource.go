package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func GetIncomeSource(c *gin.Context) {

	inComeSourceId := c.Query("id") //get countryid query from server
	//check query id is blank then display all table
	if len(inComeSourceId) == 0 {
		q := db.NewQuery("SELECT * FROM tblincomesource")
		var TblIncomeSource []struct {
			Id          int    `db:"id" json:"id"`
			Description string `db:"description" json:"description"`
		}
		q.All(&TblIncomeSource)
		c.JSON(200, TblIncomeSource)
	} else {
		//check query id is valid to convert int
		inComeSourceIdConvert, err := strconv.Atoi(inComeSourceId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found"})
			return
		}

		// check query  is in between the value // here need to check num better
		if inComeSourceIdConvert >= 1 {
			var TblIncomeSource struct {
				Id          int    `db:"id" json:"id"`
				Description string `db:"description" json:"description"`
			}

			err := db.Select().
				From("tblincomesource").
				Where(dbx.HashExp{"id": inComeSourceId}).
				One(&TblIncomeSource)

			if err != nil {
				c.JSON(404, gin.H{"error": "404 Not Found"})
			}
			c.JSON(200, TblIncomeSource)

		}

	} //end of ifelse

}
