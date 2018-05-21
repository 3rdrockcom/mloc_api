package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func GetStateId(c *gin.Context) {

	//get first query from server
	countryId := c.Query("country_id")
	stateId := c.Query("state_id")

	//require countryId for input
	if len(countryId) == 0 {
		return
	} else {
		countryIdConvert, err := strconv.Atoi(countryId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found "})
			return
		}
		if countryIdConvert < 1 || countryIdConvert > 251 { //countryId has to be valid
			return
		}

	}

	//if stateId query has no value, display all state
	if len(stateId) == 0 {

		q := db.NewQuery("Select * FROM tblstate WHERE country_id={:country_id}")
		//fetcfh all row in the struct
		var Tblstate []struct {
			StateId   int    `db:"state_id" json:"state_id"`
			State     string `db:"state" json:"state"`
			StateCode string `db:"state_code" json:"state_code"`
			CountryId int    `db:"country_id" json:"country_id"`
		}
		q.Bind(dbx.Params{"country_id": countryId})
		q.All(&Tblstate)
		c.JSON(200, Tblstate)
	} else {

		//convert stateid query to int for check valid
		stateIdConvert, err := strconv.Atoi(stateId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found "})
			return
		}

		//check stateid is valid in database table
		if stateIdConvert >= 1 {

			var Tblstate struct {
				StateId   int    `db:"state_id" json:"state_id"`
				State     string `db:"state" json:"state"`
				StateCode string `db:"state_code" json:"state_code"`
				CountryId int    `db:"country_id" json:"country_id"`
			}
			fmt.Printf("%d", len(stateId))
			// print out state relative to the country
			err := db.Select().
				From("tblstate").
				Where(dbx.HashExp{"state_id": stateId, "country_id": countryId}).
				One(&Tblstate)

			if err != nil {
				c.JSON(404, gin.H{"error": "404 Not found"})
				//fmt.Println(err)

			}

			c.JSON(200, Tblstate)

		} //end of check valid of stateid

	} //end of elseif, if state have a value

} //end of function
