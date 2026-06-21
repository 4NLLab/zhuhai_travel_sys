package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var admin models.AdminUser
	if err := database.DB.Where("username = ? AND status = ?", req.Username, "active").First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "用户名或密码错误"))
		return
	}

	if admin.PasswordHash != hashPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "用户名或密码错误"))
		return
	}

	now := time.Now()
	database.DB.Model(&admin).Update("last_login_at", now)

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"admin_id":     admin.ID,
		"username":     admin.Username,
		"display_name": admin.DisplayName,
		"role":         admin.Role,
	}))
}

func hashPassword(pwd string) string {
	h := sha256.Sum256([]byte(pwd + "zhuhai-salt"))
	return hex.EncodeToString(h[:])
}
