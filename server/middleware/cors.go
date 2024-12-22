// middleware/cors.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware cho Gin
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Thiết lập các header CORS
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Chấp nhận mọi nguồn gốc
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Nếu là yêu cầu OPTIONS, phản hồi ngay lập tức
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Chuyển tiếp đến handler tiếp theo
		c.Next()
	}
}
