package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckSubnetMiddleware(trustedSubnet *net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipStr := c.Request.Header.Get("X-Real-IP")
		if ipStr == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if !trustedSubnet.Contains(ip) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
