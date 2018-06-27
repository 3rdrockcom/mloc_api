package logger

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Tags to constructed the logger format.
		//
		// - time_unix
		// - time_unix_nano
		// - time_rfc3339
		// - time_rfc3339_nano
		// - time_custom
		// - id (Request ID)
		// - remote_ip
		// - uri
		// - host
		// - method
		// - path
		// - referer
		// - user_agent
		// - status
		// - latency (In nanoseconds)
		// - latency_human (Human readable)
		// - bytes_in (Bytes received)
		// - bytes_out (Bytes sent)
		// - header:<NAME>
		// - query:<NAME>
		// - form:<NAME>
		//
		// Example "${remote_ip} ${status}"
		//
		// Optional. Default value DefaultLoggerConfig.Format.
		Format map[string]string `yaml:"format"`

		// Optional. Default value DefaultLoggerConfig.CustomTimeFormat.
		CustomTimeFormat string `yaml:"custom_time_format"`

		// Logger is a logrus instance
		Logger *logrus.Logger
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: map[string]string{
			"time":          "time_rfc3339_nano",
			"id":            "id",
			"remote_ip":     "remote_ip",
			"host":          "host",
			"method":        "method",
			"uri":           "uri",
			"status":        "status",
			"latency":       "latency",
			"latency_human": "latency_human",
			"bytes_in":      "bytes_in",
			"bytes_out":     "bytes_out",
		},
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Format == nil {
		config.Format = DefaultLoggerConfig.Format
	}
	if config.Logger == nil {
		panic("echo: logrus middleware requires logrus instance to operate")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			fields := make(map[string]interface{})
			for field, tag := range config.Format {
				switch tag {
				case "time_unix":
					fields[field] = strconv.FormatInt(time.Now().Unix(), 10)
				case "time_unix_nano":
					fields[field] = strconv.FormatInt(time.Now().UnixNano(), 10)
				case "time_rfc3339":
					fields[field] = time.Now().Format(time.RFC3339)
				case "time_rfc3339_nano":
					fields[field] = time.Now().Format(time.RFC3339Nano)
				case "time_custom":
					fields[field] = time.Now().Format(config.CustomTimeFormat)
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					fields[field] = id
				case "remote_ip":
					fields[field] = c.RealIP()
				case "host":
					fields[field] = req.Host
				case "uri":
					fields[field] = req.RequestURI
				case "method":
					fields[field] = req.Method
				case "path":
					p := req.URL.Path
					if p == "" {
						p = "/"
					}
					fields[field] = p
				case "referer":
					fields[field] = req.Referer()
				case "user_agent":
					fields[field] = req.UserAgent()
				case "status":
					fields[field] = res.Status
				case "latency":
					l := stop.Sub(start)
					fields[field] = strconv.FormatInt(int64(l), 10)
				case "latency_human":
					fields[field] = stop.Sub(start).String()
				case "bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					fields[field] = cl
				case "bytes_out":
					fields[field] = strconv.FormatInt(res.Size, 10)
				default:
					switch {
					case strings.HasPrefix(tag, "header:"):
						fields[field] = []byte(c.Request().Header.Get(tag[7:]))
					case strings.HasPrefix(tag, "query:"):
						fields[field] = []byte(c.QueryParam(tag[6:]))
					case strings.HasPrefix(tag, "form:"):
						fields[field] = []byte(c.FormValue(tag[5:]))
					case strings.HasPrefix(tag, "cookie:"):
						cookie, err := c.Cookie(tag[7:])
						if err == nil {
							fields[field] = []byte(cookie.Value)
						}
					}

				}
			}

			config.Logger.
				WithFields(fields).
				Info("API Request")
			return
		}
	}
}
