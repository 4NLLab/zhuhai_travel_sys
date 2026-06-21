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
	var user models.User
	if err := database.DB.First(&user, currentUserID(c)).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "用户不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(user))
}

// UserRealnameSubmit 提交实名认证
func UserRealnameSubmit(c *gin.Context) {
	var req struct {
		RealName string `json:"real_name" binding:"required"`
		IdCardNo string `json:"id_card_no" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", currentUserID(c)).Updates(map[string]interface{}{
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
	var travelers []models.Traveler
	database.DB.Where("user_id = ?", currentUserID(c)).Order("is_default DESC").Find(&travelers)
	c.JSON(http.StatusOK, dto.Success(travelers))
}

// TravelerCreate 添加出游人
func TravelerCreate(c *gin.Context) {
	var t models.Traveler
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	t.UserID = currentUserID(c)
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
	if err := database.DB.Where("id = ? AND user_id = ?", id, currentUserID(c)).First(&t).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "出游人不存在"))
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	t.UserID = currentUserID(c)
	database.DB.Save(&t)
	c.JSON(http.StatusOK, dto.Success(t))
}

// TravelerDelete 删除出游人
func TravelerDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Where("user_id = ?", currentUserID(c)).Delete(&models.Traveler{}, id)
	c.JSON(http.StatusOK, dto.Success(nil))
}

// ==================== 收藏 ====================

// FavoriteList 收藏列表
func FavoriteList(c *gin.Context) {
	var favs []models.UserFavorite
	database.DB.Where("user_id = ?", currentUserID(c)).Preload("Product").Find(&favs)
	c.JSON(http.StatusOK, dto.Success(favs))
}

// FavoriteToggle 切换收藏
func FavoriteToggle(c *gin.Context) {
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	var fav models.UserFavorite
	err := database.DB.Where("user_id = ? AND product_id = ?", currentUserID(c), req.ProductID).First(&fav).Error
	if err == nil {
		database.DB.Delete(&fav)
		c.JSON(http.StatusOK, dto.Success(gin.H{"favorited": false}))
	} else {
		newFav := models.UserFavorite{UserID: currentUserID(c), ProductID: req.ProductID}
		if err := database.DB.Create(&newFav).Error; err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(400, "收藏失败"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(gin.H{"favorited": true}))
	}
}

// ==================== 发票抬头 ====================

// InvoiceTitleList 发票抬头列表
func InvoiceTitleList(c *gin.Context) {
	var titles []models.InvoiceTitle
	database.DB.Where("user_id = ?", currentUserID(c)).Find(&titles)
	c.JSON(http.StatusOK, dto.Success(titles))
}

// InvoiceTitleCreate 添加发票抬头
func InvoiceTitleCreate(c *gin.Context) {
	var t models.InvoiceTitle
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	t.UserID = currentUserID(c)
	if t.IsDefault == 1 {
		database.DB.Model(&models.InvoiceTitle{}).Where("user_id = ?", t.UserID).Update("is_default", 0)
	}
	if err := database.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "添加失败"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(t))
}

// InvoiceTitleDelete 删除发票抬头
func InvoiceTitleDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Where("user_id = ?", currentUserID(c)).Delete(&models.InvoiceTitle{}, id)
	c.JSON(http.StatusOK, dto.Success(nil))
}

// InvoiceCreate 申请开票
func InvoiceCreate(c *gin.Context) {
	var req struct {
		OrderID        uint64  `json:"order_id" binding:"required"`
		InvoiceTitleID *uint64 `json:"invoice_title_id"`
		Amount         float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", req.OrderID, currentUserID(c)).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return
	}
	if order.Status != "paid" && order.Status != "completed" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "仅已支付订单可申请开票"))
		return
	}
	if req.Amount <= 0 || req.Amount > order.PaidAmount {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "开票金额不合法"))
		return
	}
	if req.InvoiceTitleID != nil {
		var title models.InvoiceTitle
		if err := database.DB.Where("id = ? AND user_id = ?", *req.InvoiceTitleID, currentUserID(c)).First(&title).Error; err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(400, "发票抬头不存在"))
			return
		}
	}
	var invoiced float64
	database.DB.Model(&models.Invoice{}).
		Where("order_id = ? AND status <> ?", req.OrderID, "rejected").
		Select("COALESCE(SUM(amount), 0)").Scan(&invoiced)
	if invoiced+req.Amount > order.PaidAmount {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "累计开票金额超过实付金额"))
		return
	}
	inv := models.Invoice{
		UserID:         currentUserID(c),
		OrderID:        req.OrderID,
		InvoiceTitleID: req.InvoiceTitleID,
		Amount:         req.Amount,
		InvoiceNo:      strPtr(""),
	}
	if err := database.DB.Create(&inv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "申请失败"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(inv))
}

func strPtr(s string) *string { return &s }

func currentUserID(c *gin.Context) uint64 {
	if v, ok := c.Get("actor_id"); ok {
		if id, ok := v.(uint64); ok {
			return id
		}
	}
	return 0
}
