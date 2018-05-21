package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func GetPayFrequency(c *gin.Context) {
	payFrequencyId := c.Query("id")
	//check if payfrequency id query is blank,display all row in table
	if len(payFrequencyId) == 0 {
		q := db.NewQuery("SELECT * FROM tblpayfrequency")
		//fetch all row
		var TblPayFrequency []struct {
			Id          int    `db:"id" json:"id"`
			Description string `db:"description" json:"description"`
		}
		q.All(&TblPayFrequency)
		c.JSON(200, TblPayFrequency)
	} else { // contains queryid

		//check queryid is interge. if it is not, return
		payfrequencyIdConvert, err := strconv.Atoi(payFrequencyId)
		if err != nil {
			c.JSON(400, gin.H{"error": "id is not interge"})
			return
		}

		if payfrequencyIdConvert >= 1 {
			var TblPayFrequency struct {
				Id          int    `db:"id" json:"id"`
				Description string `db:"description" json:"description"`
			}

			err := db.Select("id", "description").
				From("tblpayfrequency").
				Where(dbx.HashExp{"id": payFrequencyId}).
				One(&TblPayFrequency)

			if err != nil {
				c.JSON(404, gin.H{"error": "404 Not Found"})

			}
			c.JSON(200, TblPayFrequency)
		} // end of check value between

	} //end of if else

}
