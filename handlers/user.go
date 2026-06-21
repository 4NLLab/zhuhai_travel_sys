package handlers

import (
	"net/http"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// ==================== 用户 ====================

// UserProfile 获取用户信息
func UserProfile(c *gin.Context) {
	id := c.Query("user_id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "用户不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(user))
}

// UserRealnameSubmit 提交实名认证
func UserRealnameSubmit(c *gin.Context) {
	var req struct {
		UserID   uint64 `json:"user_id" binding:"required"`
		RealName string `json:"real_name" binding:"required"`
		IdCardNo string `json:"id_card_no" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", req.UserID).Updates(map[string]interface{}{
		"real_name":       req.RealName,
		"id_card_no":      req.IdCardNo,
		"realname_status": "verified",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "提交失败"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{"realname_status": "verified"}))
}

// ==================== 出游人 ====================

// TravelerList 出游人列表
func TravelerList(c *gin.Context) {
	userID := c.Query("user_id")
	var travelers []models.Traveler
	database.DB.Where("user_id = ?", userID).Order("is_default DESC").Find(&travelers)
	c.JSON(http.StatusOK, dto.Success(travelers))
}

// TravelerCreate 添加出游人
func TravelerCreate(c *gin.Context) {
	var t models.Traveler
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if t.IsDefault == 1 {
		database.DB.Model(&models.Traveler{}).Where("user_id = ?", t.UserID).Update("is_default", 0)
	}
	if err := database.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "添加失败"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(t))
}

// TravelerUpdate 编辑出游人
func TravelerUpdate(c *gin.Context) {
	id := c.Param("id")
	var t models.Traveler
	if err := database.DB.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "出游人不存在"))
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	database.DB.Save(&t)
	c.JSON(http.StatusOK, dto.Success(t))
}

// TravelerDelete 删除出游人
func TravelerDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Traveler{}, id)
	c.JSON(http.StatusOK, dto.Success(nil))
}

// ==================== 收藏 ====================

// FavoriteList 收藏列表
func FavoriteList(c *gin.Context) {
	userID := c.Query("user_id")
	var favs []models.UserFavorite
	database.DB.Where("user_id = ?", userID).Preload("Product").Find(&favs)
	c.JSON(http.StatusOK, dto.Success(favs))
}

// FavoriteToggle 切换收藏
func FavoriteToggle(c *gin.Context) {
	var req struct {
		UserID    uint64 `json:"user_id" binding:"required"`
		ProductID uint64 `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	var fav models.UserFavorite
	err := database.DB.Where("user_id = ? AND product_id = ?", req.UserID, req.ProductID).First(&fav).Error
	if err == nil {
		database.DB.Delete(&fav)
		c.JSON(http.StatusOK, dto.Success(gin.H{"favorited": false}))
	} else {
		newFav := models.UserFavorite{UserID: req.UserID, ProductID: req.ProductID}
		database.DB.Create(&newFav)
		c.JSON(http.StatusOK, dto.Success(gin.H{"favorited": true}))
	}
}

// ==================== 发票抬头 ====================

// InvoiceTitleList 发票抬头列表
func InvoiceTitleList(c *gin.Context) {
	userID := c.Query("user_id")
	var titles []models.InvoiceTitle
	database.DB.Where("user_id = ?", userID).Find(&titles)
	c.JSON(http.StatusOK, dto.Success(titles))
}

// InvoiceTitleCreate 添加发票抬头
func InvoiceTitleCreate(c *gin.Context) {
	var t models.InvoiceTitle
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if t.IsDefault == 1 {
		database.DB.Model(&models.InvoiceTitle{}).Where("user_id = ?", t.UserID).Update("is_default", 0)
	}
	database.DB.Create(&t)
	c.JSON(http.StatusOK, dto.Success(t))
}

// InvoiceTitleDelete 删除发票抬头
func InvoiceTitleDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.InvoiceTitle{}, id)
	c.JSON(http.StatusOK, dto.Success(nil))
}

// InvoiceCreate 申请开票
func InvoiceCreate(c *gin.Context) {
	var req struct {
		UserID         uint64  `json:"user_id" binding:"required"`
		OrderID        uint64  `json:"order_id" binding:"required"`
		InvoiceTitleID *uint64 `json:"invoice_title_id"`
		Amount         float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	inv := models.Invoice{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		InvoiceTitleID: req.InvoiceTitleID,
		Amount:  req.Amount,
		InvoiceNo: strPtr(""),
	}
	database.DB.Create(&inv)
	c.JSON(http.StatusOK, dto.Success(inv))
}

func strPtr(s string) *string { return &s }
