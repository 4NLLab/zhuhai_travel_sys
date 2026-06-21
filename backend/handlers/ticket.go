package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// TicketDetail 票详情
func TicketDetail(c *gin.Context) {
	id := c.Param("id")
	var ticket models.Ticket
	if err := database.DB.Preload("OrderItem").First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "票不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(ticket))
}

// TicketByQR 通过二维码哈希获取票信息（核销前查询）
func TicketByQR(c *gin.Context) {
	hash := c.Query("qr_hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少二维码哈希"))
		return
	}
	var ticket models.Ticket
	if err := database.DB.Preload("OrderItem").
		Where("qr_token_hash = ?", hash).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "票不存在"))
		return
	}
	if ticket.Status != "valid" {
		c.JSON(http.StatusOK, dto.Fail(409, "票已失效，状态: "+ticket.Status))
		return
	}
	c.JSON(http.StatusOK, dto.Success(ticket))
}

// TicketVerify 核销电子票
func TicketVerify(c *gin.Context) {
	var req struct {
		TicketID         uint64  `json:"ticket_id" binding:"required"`
		VerifierAdminID  *uint64 `json:"verifier_admin_id"`
		VerifyLocation   string  `json:"verify_location"`
		VerifyNote       string  `json:"verify_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var ticket models.Ticket
	if err := database.DB.First(&ticket, req.TicketID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "票不存在"))
		return
	}
	if ticket.Status != "valid" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "票已失效"))
		return
	}

	now := timeNow()
	tx := database.DB.Begin()

	// 更新票状态
	var loc *string
	if req.VerifyLocation != "" {
		loc = &req.VerifyLocation
	}
	if err := tx.Model(&ticket).Updates(map[string]interface{}{
		"status":  "used",
		"used_at": now,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "核销失败"))
		return
	}

	// 写核销记录
	var note *string
	if req.VerifyNote != "" {
		note = &req.VerifyNote
	}
	verification := models.TicketVerification{
		TicketID:        ticket.ID,
		VerifierAdminID: req.VerifierAdminID,
		VerifyLocation:  loc,
		VerifyResult:    "success",
		VerifyNote:      note,
	}
	tx.Create(&verification)

	// 生成司机佣金 — 根据该票的订单是否有推广司机（核销才算成功）
	var item models.OrderItem
	if err := tx.First(&item, ticket.OrderItemID).Error; err == nil {
		var order models.Order
		if err := tx.First(&order, item.OrderID).Error; err == nil && order.DriverID != nil {
			var count int64
			tx.Model(&models.DriverCommission{}).
				Where("order_item_id = ? AND driver_id = ?", item.ID, *order.DriverID).
				Count(&count)
			if count == 0 {
				var driver models.Driver
				if tx.First(&driver, *order.DriverID).Error == nil {
					commAmount := float64(item.Quantity) * item.UnitPrice * driver.CommissionRate
					if commAmount > 0 {
						commNo := fmt.Sprintf("CM%s%04d", timeNow().Format("20060102150405"), rand.Intn(10000))
						tx.Create(&models.DriverCommission{
							DriverID: *order.DriverID, OrderID: order.ID, OrderItemID: item.ID,
							CommissionNo: commNo,
							BaseAmount:   float64(item.Quantity) * item.UnitPrice,
							Rate:         driver.CommissionRate,
							CommissionAmount: commAmount,
						})
					}
				}
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"ticket_id":    ticket.ID,
		"ticket_no":    ticket.TicketNo,
		"status":       "used",
		"verified_at":  now,
	}))
}

// VerificationHistory 核销记录查询
func VerificationHistory(c *gin.Context) {
	ticketID := c.Query("ticket_id")
	var list []models.TicketVerification
	database.DB.Where("ticket_id = ?", ticketID).Order("created_at DESC").Find(&list)
	c.JSON(http.StatusOK, dto.Success(list))
}

