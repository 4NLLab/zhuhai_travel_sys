package handlers

import (
	"net/http"
	"strconv"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"
	"zhuhai_travel_backend/security"

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

	if !security.VerifyPassword(req.Password, admin.PasswordHash) {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "用户名或密码错误"))
		return
	}

	now := time.Now()
	database.DB.Model(&admin).Update("last_login_at", now)

	token, err := security.GenerateToken(config.Load().JWTSecret, admin.ID, "admin", admin.DisplayName, tokenTTL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "生成登录凭证失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"admin_id":     admin.ID,
		"username":     admin.Username,
		"display_name": admin.DisplayName,
		"role":         admin.Role,
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(tokenTTL().Seconds()),
	}))
}

// AdminProfile 返回当前超级管理员登录态，用于前端刷新后的 token 校验。
func AdminProfile(c *gin.Context) {
	adminID := currentActorID(c)
	var admin models.AdminUser
	if err := database.DB.Where("id = ? AND status = ?", adminID, "active").First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "管理员账号已停用或不存在"))
		return
	}
	if admin.Role != "super_admin" {
		c.JSON(http.StatusForbidden, dto.Fail(403, "仅超级管理员可访问"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"admin_id":      admin.ID,
		"username":      admin.Username,
		"display_name":  admin.DisplayName,
		"role":          admin.Role,
		"last_login_at": admin.LastLoginAt,
	}))
}

func UserPhoneLogin(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var user models.User
	err := database.DB.Where("phone = ? AND status = ?", req.Phone, "active").First(&user).Error
	if err != nil {
		user = models.User{
			Phone:    &req.Phone,
			Nickname: strPtr(req.Nickname),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建用户失败"))
			return
		}
	}

	now := time.Now()
	database.DB.Model(&user).Update("last_login_at", now)
	name := req.Phone
	if user.Nickname != nil && *user.Nickname != "" {
		name = *user.Nickname
	}
	token, err := security.GenerateToken(config.Load().JWTSecret, user.ID, "user", name, tokenTTL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "生成登录凭证失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"user_id":      user.ID,
		"phone":        security.MaskPhone(req.Phone),
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(tokenTTL().Seconds()),
	}))
}

func tokenTTL() time.Duration {
	hours, err := strconv.Atoi(config.Load().TokenTTLHours)
	if err != nil || hours <= 0 {
		hours = 168
	}
	return time.Duration(hours) * time.Hour
}
