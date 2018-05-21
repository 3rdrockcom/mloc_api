package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	//	"net/http"
)

func GetCountryId(c *gin.Context) {
	// need to check the login authelization and head api key

	countryId := c.Query("country_id") //get countryid query from server

	//if countryid is empty string
	if len(countryId) == 0 {
		fmt.Println(len(countryId))
		q := db.NewQuery("SELECT * FROM tblcountry")
		var Tblcountry []struct {
			CountryId        int    `db:"country_id" json:"country_id"`
			Name             string `db:"name" json:"name"`
			IsoCode2         string `db:"iso_code_2" json:"iso_code_2"`
			IsoCode3         string `db:"iso_code_3" json:"iso_coe_3"`
			AddressFormat    string `db:address_format" json:"address_format"`
			PostCodeRequired int    `db:"postcode_required" json:"postcode_required"`
			Status           int    `db:"status" json:"status"`
			MobilePrefix     string `db:"mobile_prefix" json:"mobile_prefix"`
		}

		q.All(&Tblcountry)
		c.JSON(200, Tblcountry)

	} else {
		//countryId has a value
		//convert countryid  value to interge to check is valid or not
		countryIdConvert, err := strconv.Atoi(countryId)
		if err != nil {
			c.JSON(400, gin.H{"error": "Not Found "})
			return
		}

		//if integer countryid is between 1 to 251, then it is valid in database
		if countryIdConvert >= 1 {
			var Tblcountry struct {
				CountryId        int    `db:"country_id" json:"country_id"`
				Name             string `db:"name" json:"name"`
				IsoCode2         string `db:"iso_code_2" json:"iso_code_2"`
				IsoCode3         string `db:"iso_code_3" json:"iso_coe_3"`
				AddressFormat    string `db:address_format" json:"address_format"`
				PostCodeRequired int    `db:"postcode_required" json:"postcode_required"`
				Status           int    `db:"status" json:"status"`
				MobilePrefix     string `db:"mobile_prefix" json:"mobile_prefix"`
			}

			err := db.Select("country_id", "name", "iso_code_2", "iso_code_3", "address_format", "postcode_required", "status", "mobile_prefix").
				From("tblcountry").
				Where(dbx.HashExp{"country_id": countryId}).
				One(&Tblcountry)

			if err != nil {
				c.JSON(404, gin.H{"error": "404 Not Found"})
			}
			c.JSON(200, Tblcountry)
		} // end of if countryid is between 1 to 251
	} //end of else of input countryid value

} // end of function
