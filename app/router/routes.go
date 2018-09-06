package router

// appendRoutes registers routes in the router
func (r *Router) appendRoutes() {
	// API
	api := r.e.Group("/api")
	api.Use(r.mwBasicAuth())

	// API, Version 1
	v1 := api.Group("/v1")

	// Endpoints for auth
	v1.GET("/login/get_customer_key", r.c.GetCustomerKey, r.mwKeyAuth("login", "cust_unique_id"))
	v1.GET("/config/generate_customer_key", r.c.GenerateCustomerKey, r.mwKeyAuth("registration", "cust_unique_id"))

	// Endpoints for customer
	v1.GET("/customer/get_customer", r.c.GetCustomer, r.mwKeyAuth("default", "cust_unique_id"))
	v1.GET("/customer/get_transaction_history", r.c.GetTransactionHistory, r.mwKeyAuth("default", "R1"))
	v1.GET("/customer/get_customer_loan", r.c.GetCustomerLoan, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/customer_basic", r.c.PostCustomerBasic, r.mwKeyAuth("default", "R16"))
	v1.POST("/customer/customer_additional", r.c.PostCustomerAdditional, r.mwKeyAuth("default", "R8"))
	v1.POST("/customer/accept_terms_and_condition", r.c.PostAcceptTermsAndConditions, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/credit_line_application", r.c.PostCreditLineApplication, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/compute_loan", r.c.PostComputeLoan, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/loan_application", r.c.PostLoanApplication, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/pay_loan", r.c.PostPayLoan, r.mwKeyAuth("default", "R1"))

	v1.GET("/customer/get_bank_account", r.c.GetBankAccount, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/create_bank_account", r.c.PostCreateBankAccount, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/update_bank_account", r.c.PostUpdateBankAccount, r.mwKeyAuth("default", "R1"))
	v1.POST("/customer/delete_bank_account", r.c.PostDeleteBankAccount, r.mwKeyAuth("default", "R1"))
	v1.GET("/customer/get_bank_institution", r.c.GetBankInstitutions, r.mwKeyAuth("default", "R1"))

	// Endpoints for payment
	v1.GET("/payment/get_disbursement_method", r.c.GetDisbursementMethod, r.mwKeyAuth("default", "R1"))
	v1.GET("/payment/get_collection_method", r.c.GetCollectionMethod, r.mwKeyAuth("default", "R1"))

	// Endpoints for lookup
	v1.GET("/lookup/get_country", r.c.GetCountry, r.mwKeyAuth("default", ""))
	v1.GET("/lookup/get_income_source", r.c.GetIncomeSource, r.mwKeyAuth("default", ""))
	v1.GET("/lookup/get_pay_frequency", r.c.GetPayFrequency, r.mwKeyAuth("default", ""))
	v1.GET("/lookup/get_state", r.c.GetState, r.mwKeyAuth("default", ""))
	v1.GET("/lookup/get_city", r.c.GetCity, r.mwKeyAuth("default", ""))

}
