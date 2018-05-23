package router

import (
	"net"
	"os"

	"github.com/epointpayment/mloc_api_go/app/controllers"

	"github.com/labstack/echo"
)

// Router manages the applications routing functions
type Router struct {
	c *controllers.Controllers
	e *echo.Echo
}

// NewRouter creates an instance of the service
func NewRouter(c *controllers.Controllers) *Router {
	r := &Router{}

	r.e = echo.New()
	r.c = c

	r.appendMiddleware()
	r.appendRoutes()
	r.appendErrorHandler()

	return r
}

func (r *Router) Run() error {
	// Get host information
	host := ""
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}

	//Get port information
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Create an address for the router to use
	address := net.JoinHostPort(host, port)

	// Start routing service
	return r.e.Start(address)
}
