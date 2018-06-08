package router

import (
	"net"
	"strconv"

	"github.com/epointpayment/mloc_api_go/app/config"
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

	r.c = c

	// Initialize router
	r.e = echo.New()
	r.e.HideBanner = true

	r.appendMiddleware()
	r.appendRoutes()
	r.appendErrorHandler()

	return r
}

func (r *Router) Run() error {
	// Get config information
	host := config.Get().Server.Host
	port := strconv.FormatInt(config.Get().Server.Port, 10)

	// Create an address for the router to use
	address := net.JoinHostPort(host, port)

	// Start routing service
	return r.e.Start(address)
}
