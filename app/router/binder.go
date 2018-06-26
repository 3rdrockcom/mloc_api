package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/monoculum/formam"
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {

		req := c.Request()
		ctype := req.Header.Get(echo.HeaderContentType)

		switch {
		case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm):
			params, err := c.FormParams()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			dec := formam.NewDecoder(&formam.DecoderOptions{TagName: "form", IgnoreUnknownKeys: true})
			if err := dec.Decode(params, i); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return nil
		}

		return
	}

	// Define your custom implementation

	return
}
