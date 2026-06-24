package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"
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

func AdminRoleRequired(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = true
	}

	return func(c *gin.Context) {
		if actorRole, _ := c.Get("actor_role"); actorRole != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail(403, "无权访问该资源"))
			return
		}

		actorID, ok := actorIDFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "登录已过期或无效"))
			return
		}

		var admin models.AdminUser
		if err := database.DB.Where("id = ? AND status = ?", actorID, "active").First(&admin).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "管理员账号已停用或不存在"))
			return
		}

		if len(roleSet) > 0 && !roleSet[admin.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail(403, "仅超级管理员可访问"))
			return
		}

		c.Set("actor_admin_role", admin.Role)
		c.Set("actor_name", admin.DisplayName)
		c.Next()
	}
}

func DriverStatusRequired(allowedStatuses ...string) gin.HandlerFunc {
	statusSet := make(map[string]bool, len(allowedStatuses))
	for _, status := range allowedStatuses {
		statusSet[status] = true
	}

	return func(c *gin.Context) {
		if actorRole, _ := c.Get("actor_role"); actorRole != "driver" {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail(403, "无权访问该资源"))
			return
		}

		actorID, ok := actorIDFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "登录已过期或无效"))
			return
		}

		var driver models.Driver
		if err := database.DB.First(&driver, actorID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(401, "司机账号不存在"))
			return
		}

		if len(statusSet) > 0 && !statusSet[driver.Status] {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.Fail(403, "司机账号未启用"))
			return
		}

		c.Set("actor_name", driver.Name)
		c.Set("actor_driver_status", driver.Status)
		c.Next()
	}
}

func actorIDFromContext(c *gin.Context) (uint64, bool) {
	id, ok := c.Get("actor_id")
	if !ok {
		return 0, false
	}
	typed, ok := id.(uint64)
	return typed, ok && typed > 0
}
