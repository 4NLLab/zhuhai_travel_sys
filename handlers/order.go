package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建订单（含库存锁定事务）
func CreateOrder(c *gin.Context) {
	var req struct {
		UserID       uint64 `json:"user_id" binding:"required"`
		SkuID        uint64 `json:"sku_id" binding:"required"`
		ScheduleID   uint64 `json:"schedule_id"`
		Quantity     uint   `json:"quantity" binding:"required,min=1"`
		DriverQRCode string `json:"driver_qr_code"`
		Travelers    []struct {
			Name  string `json:"name" binding:"required"`
			Phone string `json:"phone"`
			IdNo  string `json:"id_no" binding:"required"`
		} `json:"travelers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误: "+err.Error()))
		return
	}

	var sku models.ProductSKU
	if err := database.DB.First(&sku, req.SkuID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "SKU 不存在"))
		return
	}
	if sku.Status != "active" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "SKU 已下架"))
		return
	}

	var schedule models.ProductSchedule
	if err := database.DB.First(&schedule, req.ScheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "排期不存在"))
		return
	}
	avail := schedule.TotalStock - schedule.LockedStock - schedule.SoldStock
	if int(req.Quantity) > avail {
		c.JSON(http.StatusBadRequest, dto.Fail(400,
			fmt.Sprintf("库存不足，剩余: %d", avail)))
		return
	}

	var product models.Product
	database.DB.First(&product, sku.ProductID)

	tx := database.DB.Begin()

	// 行锁锁定库存 — 用条件更新防止并发超卖
	result := tx.Model(&models.ProductSchedule{}).
		Where("id = ? AND total_stock - locked_stock - sold_stock >= ?", schedule.ID, req.Quantity).
		UpdateColumn("locked_stock", schedule.LockedStock+int(req.Quantity))
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, dto.Fail(400, "库存已被抢光"))
		return
	}

	orderNo := generateOrderNo()
	totalPrice := float64(req.Quantity) * sku.SalePrice

	// 检查司机扫码
	var driverID, driverQRCodeID *uint64
	source := "miniapp"
	if req.DriverQRCode != "" {
		var qr models.DriverQRCode
		if err := database.DB.Where("code = ? AND status = ?", req.DriverQRCode, "active").First(&qr).Error; err == nil {
			driverID = &qr.DriverID
			driverQRCodeID = &qr.ID
			source = "driver_qr"
		}
	}

	order := models.Order{
		OrderNo:        orderNo,
		UserID:         req.UserID,
		Source:         source,
		DriverID:       driverID,
		DriverQRCodeID: driverQRCodeID,
		TotalAmount:    totalPrice,
		PayableAmount:  totalPrice,
		ContactName:    &req.Travelers[0].Name,
		ContactPhone:   &req.Travelers[0].Phone,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建订单失败"))
		return
	}

	var travelDate string
	if schedule.TravelDate != "" {
		travelDate = schedule.TravelDate
	}
	item := models.OrderItem{
		OrderID:      order.ID,
		ProductID:    sku.ProductID,
		SkuID:        sku.ID,
		ScheduleID:   &schedule.ID,
		ProductTitle: product.Title,
		SkuName:      sku.SkuName,
		TravelDate:   &travelDate,
		StartTime:    schedule.StartTime,
		Quantity:     req.Quantity,
		UnitPrice:    sku.SalePrice,
		TotalPrice:   totalPrice,
	}
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建订单项失败"))
		return
	}

	for _, t := range req.Travelers {
		ot := models.OrderTraveler{
			OrderItemID: item.ID,
			Name:        t.Name,
			Phone:       &t.Phone,
			IdNo:        t.IdNo,
		}
		tx.Create(&ot)
	}

	// 预创建支付记录
	paymentNo := "PY" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", rand.Intn(10000))
	payment := models.Payment{
		OrderID:   order.ID,
		PaymentNo: paymentNo,
		Channel:   "wechat",
		Amount:    totalPrice,
	}
	tx.Create(&payment)

	tx.Commit()

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"order_no":    orderNo,
		"payment_no":  paymentNo,
		"amount":      totalPrice,
		"order_id":    order.ID,
	}))
}

// PaymentCallback 支付回调（微信/支付宝异步通知）
func PaymentCallback(c *gin.Context) {
	var req struct {
		PaymentNo     string  `json:"payment_no" binding:"required"`
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount" binding:"required"`
		RawPayload    string  `json:"raw_payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var payment models.Payment
	if err := database.DB.Where("payment_no = ? AND status = ?", req.PaymentNo, "pending").
		First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "支付单不存在或已处理"))
		return
	}

	now := timeNow()
	tx := database.DB.Begin()

	// 更新支付记录
	tx.Model(&payment).Updates(map[string]interface{}{
		"status":         "success",
		"transaction_id": req.TransactionID,
		"paid_at":        now,
		"raw_payload":    req.RawPayload,
	})

	// 更新订单状态
	tx.Model(&models.Order{}).Where("id = ?", payment.OrderID).Updates(map[string]interface{}{
		"status":       "paid",
		"paid_amount":  req.Amount,
		"paid_at":      now,
	})

	// 库存：locked → sold
	var order models.Order
	tx.First(&order, payment.OrderID)
	var items []models.OrderItem
	tx.Where("order_id = ?", order.ID).Find(&items)
	for _, item := range items {
		if item.ScheduleID != nil {
			tx.Model(&models.ProductSchedule{}).
				Where("id = ?", *item.ScheduleID).
				Updates(map[string]interface{}{
					"locked_stock": database.DB.Raw("locked_stock - ?", item.Quantity),
					"sold_stock":   database.DB.Raw("sold_stock + ?", item.Quantity),
				})
		}
		// 生成电子票
		for i := 0; i < int(item.Quantity); i++ {
			ticketNo := "TK" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", rand.Intn(10000))
			qrpHash := "QRH" + fmt.Sprintf("%x", time.Now().UnixNano())[:16]
			tx.Create(&models.Ticket{
				OrderItemID: item.ID,
				TicketNo:    ticketNo,
				QRTokenHash: qrpHash,
			})
		}
	}

	// 司机佣金不再在支付时生成 — 改为核销时生成（见 ticket.go TicketVerify）

	tx.Commit()

	c.JSON(http.StatusOK, dto.Success(gin.H{"status": "success"}))
}

// OrderList 用户订单列表
func OrderList(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var orders []models.Order
	q.Preload("Items").Preload("Items.Travelers").Preload("Payment").
		Order("created_at DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&orders)

	c.JSON(http.StatusOK, dto.Page(orders, total, page, size))
}

// OrderDetail 订单详情
func OrderDetail(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("Items").Preload("Items.Travelers").
		Preload("Items.Tickets").Preload("Payment").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(order))
}

// ==================== helpers ====================

func generateOrderNo() string {
	return fmt.Sprintf("ZH%s%04d",
		time.Now().Format("0102150405"),
		rand.Intn(10000))
}

func timeNow() time.Time {
	return time.Now()
}
