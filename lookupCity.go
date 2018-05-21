package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func GetCity(c *gin.Context) {
	stateCode := c.Query("state_code")
	cityId := c.Query("city_id")
	fmt.Println(stateCode)
	fmt.Println(cityId)
	//require stateCode in the query
	if len(stateCode) == 0 {
		return
	} else {
		stateCodeConvert, err := strconv.Atoi(stateCode)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found"})
			return
		}
		if stateCodeConvert < 1 || stateCodeConvert > 51 {
			return
		}
	} // end of require stateCode in the query

	if len(cityId) == 0 {
		q := db.NewQuery("Select * From tblcity")

		//fetch all row in the struct
		var Tblcity []struct {
			CityId    int    `db:"city_id" json:"city_id"`
			City      string `db:"city" json:"city"`
			StateCode string `db:"state_code" json:"state_code"`
		}

		q.All(&Tblcity)
		c.JSON(200, Tblcity)

	} else {

		cityIdConvert, err := strconv.Atoi(cityId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found"})
			return
		}

		if cityIdConvert >= 1 {

			var Tblcity struct {
				CityId    int    `db:"city_id" json:"city_id"`
				City      string `db:"city" json:"city"`
				StateCode string `db:"state_code" json:"state_code"`
			}

			err := db.Select().
				From("tblcity").
				Where(dbx.HashExp{"city_id": cityId}).
				One(&Tblcity)
			if err != nil {
				c.JSON(404, gin.H{"error": "404 not found"})
			}
			//	continue here
			c.JSON(200, Tblcity)

		}
	}

}
