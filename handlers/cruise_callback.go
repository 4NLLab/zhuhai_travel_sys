package handlers

import (
	"encoding/json"
	"net/http"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/models"
	"zhuhai_travel_backend/services/jzg"

	"github.com/gin-gonic/gin"
)

func CruiseCallback(eventType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"IsSuccess": false, "Message": "JSON格式错误"})
			return
		}

		cfg := config.Load()
		verified := jzg.VerifySignature(payload, cfg.JZGAccessToken)
		orderNo := extractJZGOrderNo(payload)
		payloadBytes, _ := json.Marshal(payload)
		ip := c.ClientIP()
		logRow := models.JZGCallbackLog{
			EventType: eventType,
			OrderNo:   orderNo,
			Verified:  verified,
			RequestIP: &ip,
			Payload:   string(payloadBytes),
		}
		_ = database.DB.Create(&logRow).Error

		if !verified {
			c.JSON(http.StatusUnauthorized, gin.H{"IsSuccess": false, "Message": "签名验证失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"IsSuccess": true, "Message": "接收成功"})
	}
}

func extractJZGOrderNo(payload map[string]interface{}) *string {
	if value, ok := payload["orderNo"]; ok {
		orderNo := stringFromAny(value)
		if orderNo != "" {
			return &orderNo
		}
	}
	if orderInfo, ok := payload["orderInfo"].(map[string]interface{}); ok {
		orderNo := stringFromAny(orderInfo["orderNo"])
		if orderNo != "" {
			return &orderNo
		}
	}
	return nil
}

func stringFromAny(value interface{}) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	b, _ := json.Marshal(value)
	return string(b)
}
