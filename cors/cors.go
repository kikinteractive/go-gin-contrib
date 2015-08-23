// Forked and adapted from https://github.com/tommy351/gin-cors

package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	defaultAllowHeaders = []string{"Origin", "Accept", "Accept-Encoding", "Content-Type", "Content-Length", "Authorization"}
	defaultAllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}
)

// Options stores configurations
type Options struct {
	AllowOrigins     []string
	AllowCredentials bool
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           time.Duration
}

// Middleware sets CORS headers for every request.
// NOTE: by default (if no specific options are specified), it is permissive.
func Middleware(options Options) gin.HandlerFunc {
	if options.AllowHeaders == nil {
		options.AllowHeaders = defaultAllowHeaders
	}

	if options.AllowMethods == nil {
		options.AllowMethods = defaultAllowMethods
	}

	return func(c *gin.Context) {
		req := c.Request
		res := c.Writer

		// CORS headers are added whenever the browser request includes an "Origin" header
		origin := req.Header.Get("Origin")
		if origin == "" {
			return // or c.Next() ?
		}

		if len(options.AllowOrigins) > 0 {
			// NOTE: The string "*" cannot be used for a resource that supports credentials
			res.Header().Set("Access-Control-Allow-Origin", strings.Join(options.AllowOrigins, " "))
		} else {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if options.AllowCredentials {
			res.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if len(options.ExposeHeaders) > 0 {
			res.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposeHeaders, ","))
		}

		// Handle preflight request if applicable.
		if req.Method == "OPTIONS" {
			preflightRequestMethod := req.Header.Get("Access-Control-Request-Method")
			preflightRequestHeaders := req.Header.Get("Access-Control-Request-Headers")

			if len(options.AllowMethods) > 0 {
				res.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowMethods, ","))
			} else if preflightRequestMethod != "" {
				res.Header().Set("Access-Control-Allow-Methods", preflightRequestMethod)
			}

			if len(options.AllowHeaders) > 0 {
				res.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowHeaders, ","))
			} else if preflightRequestHeaders != "" {
				res.Header().Set("Access-Control-Allow-Headers", preflightRequestHeaders)
			}

			if options.MaxAge > time.Duration(0) {
				res.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(options.MaxAge/time.Second), 10))
			}

			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
