package main

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func CheckValidCustUniqueId(queryValue string) bool {

	var GetValidQuery struct {
		value string
	}
	err := db.Select("cust_unique_id").
		From("tblcustomerbasicinfo").
		Where(dbx.HashExp{"cust_unique_id": queryValue}).
		One(&GetValidQuery)

	if err != nil {
		//c.JSON(404, gin.H{"error": "404 key not found"})
		return false
	} else {
		//	fmt.Println("apikey login success")
		//	r.GET("/api/v1/lookup/get_country/", GetCountryId)
		return true
	}

}
