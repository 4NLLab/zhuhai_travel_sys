package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
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

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := security.ParseToken(config.Load().JWTSecret, token)
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

func AdminAuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		if c.Request.Method == http.MethodGet && strings.HasPrefix(c.Request.URL.Path, "/api/v1/admin/audit-logs") {
			c.Next()
			return
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		actorID := actorIDPtr(c)
		actorName := strPtrFromContext(c, "actor_name")
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		targetType, targetID := inferTarget(c.Request.URL.Path, bodyBytes)
		payload := buildAuditPayload(c, bodyBytes)
		payloadJSON, _ := json.Marshal(payload)
		payloadText := string(payloadJSON)

		logRow := models.AuditLog{
			ActorType:  "admin",
			ActorID:    actorID,
			ActorName:  actorName,
			Action:     auditAction(c.Request.Method, c.Request.URL.Path),
			TargetType: targetType,
			TargetID:   targetID,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			IP:         &ip,
			UserAgent:  &ua,
			Payload:    &payloadText,
		}
		if err := database.DB.Create(&logRow).Error; err != nil {
			log.Printf("audit log write failed: %v", err)
		}
	}
}

func actorIDPtr(c *gin.Context) *uint64 {
	if v, ok := c.Get("actor_id"); ok {
		if id, ok := v.(uint64); ok {
			return &id
		}
	}
	return nil
}

func strPtrFromContext(c *gin.Context, key string) *string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok && s != "" {
			return &s
		}
	}
	return nil
}

func auditAction(method, path string) string {
	switch {
	case path == "/api/v1/admin/drivers/review":
		return "admin.driver.review"
	case path == "/api/v1/admin/orders/refund":
		return "admin.order.refund"
	case path == "/api/v1/admin/withdrawals/process":
		return "admin.withdrawal.process"
	case path == "/api/v1/admin/commissions/settle":
		return "admin.commission.settle"
	case strings.HasPrefix(path, "/api/v1/admin/banners"):
		return "admin.banner." + strings.ToLower(method)
	case path == "/api/v1/tickets/verify":
		return "admin.ticket.verify"
	default:
		return "admin." + strings.ToLower(method)
	}
}

func inferTarget(path string, body []byte) (*string, *uint64) {
	var targetType string
	var id uint64
	switch {
	case strings.Contains(path, "/drivers/review"):
		targetType = "driver"
		id = uintFromJSON(body, "driver_id")
	case strings.Contains(path, "/orders/refund"):
		targetType = "order"
		id = uintFromJSON(body, "order_id")
	case strings.Contains(path, "/withdrawals/process"):
		targetType = "driver_withdrawal"
		id = uintFromJSON(body, "withdrawal_id")
	case strings.Contains(path, "/tickets/verify"):
		targetType = "ticket"
		id = uintFromJSON(body, "ticket_id")
	default:
		return nil, nil
	}
	if id == 0 {
		return &targetType, nil
	}
	return &targetType, &id
}

func uintFromJSON(body []byte, key string) uint64 {
	var data map[string]interface{}
	if len(body) == 0 || json.Unmarshal(body, &data) != nil {
		return 0
	}
	switch v := data[key].(type) {
	case float64:
		return uint64(v)
	case string:
		n, _ := strconv.ParseUint(v, 10, 64)
		return n
	default:
		return 0
	}
}

func buildAuditPayload(c *gin.Context, body []byte) map[string]interface{} {
	payload := map[string]interface{}{
		"query":       c.Request.URL.Query(),
		"status_code": c.Writer.Status(),
	}
	if len(body) > 0 {
		var data interface{}
		if json.Unmarshal(body, &data) == nil {
			payload["body"] = sanitizeAuditValue(data)
		} else {
			payload["body"] = "[non-json body]"
		}
	}
	return payload
}

func sanitizeAuditValue(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(x))
		for key, value := range x {
			lower := strings.ToLower(key)
			if strings.Contains(lower, "password") || strings.Contains(lower, "token") || strings.Contains(lower, "secret") {
				out[key] = "[已脱敏]"
				continue
			}
			if strings.Contains(lower, "url") && strings.Contains(lower, "license") {
				out[key] = "[证照图片URL已脱敏]"
				continue
			}
			out[key] = sanitizeAuditValue(value)
		}
		return out
	case []interface{}:
		out := make([]interface{}, 0, len(x))
		for _, item := range x {
			out = append(out, sanitizeAuditValue(item))
		}
		return out
	default:
		return v
	}
}
