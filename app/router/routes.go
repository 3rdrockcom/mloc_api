package router

// appendRoutes registers routes in the router
func (r *Router) appendRoutes() {
	// API
	api := r.e.Group("/api")
	api.Use(r.mwBasicAuth())

	// API, Version 1
	v1 := api.Group("/v1")

	// Endpoints for auth
	v1.GET("/login/get_customer_key", r.c.GetCustomerKey, r.mwKeyAuth("login"))
	v1.GET("/config/generate_customer_key", r.c.GenerateCustomerKey, r.mwKeyAuth("registration"))

	// Endpoints for customer
	v1.GET("/customer/get_customer", r.c.GetCustomer, r.mwKeyAuth("default"))
	// v1.POST("/customer/customer_basic", r.c.PostAddCustomer, r.mwKeyAuth("default"))
	v1.POST("/customer/customer_basic", r.c.PostCustomerBasic, r.mwKeyAuth("default"))

	// Endpoints for lookup
	v1.GET("/lookup/get_country", r.c.GetCountry, r.mwKeyAuth("default"))

}
