package models

import (
	"gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// CustomerAgreement contains information of customer id
type CustomerAgreement struct {
	ID                 int
	CustomerID         null.Int `db:"fk_customer_id"`
	MLOCAccess         zero.Int `db:"mloc_access" json:"mloc_access"`
	Registration       zero.Int `db:"registration" json:"registration"`
	TermsAndConditions zero.Int `db:"term_and_condition" json:"term_and_condition"`
}

// TableName gets customerbasicinfo table from database
func (a CustomerAgreement) TableName() string {
	return "tblCustomerAgreement"
}
