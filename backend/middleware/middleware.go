package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/security"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		latency := time.Since(start)
		log.Printf("[%d] %s %s %v", c.Writer.Status(), c.Request.Method, path, latency)
	}
}

func CORS() gin.HandlerFunc {
	cfg := config.Load()
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cfg.CORSAllowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func AuthRequired(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = true
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "请先登录"))
			return
		}

		cfg := config.Load()
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := security.ParseToken(cfg.JWTSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "登录已过期或无效"))
			return
		}
		if len(roleSet) > 0 && !roleSet[claims.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail(403, "无权访问该资源"))
			return
		}

		c.Set("actor_id", claims.SubjectID)
		c.Set("actor_role", claims.Role)
		c.Set("actor_name", claims.Name)
		c.Next()
	}
}
