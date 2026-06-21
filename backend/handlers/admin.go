package handlers

import (
	"fmt"
	"net/http"
	"time"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// ==================== 运营看板 ====================

// AdminDashboard 运营看板数据
func AdminDashboard(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	// 今日成功交易数
	var orderCount int64
	database.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = ?", "paid", today).
		Count(&orderCount)

	// 今日交易金额
	var totalAmount float64
	database.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = ?", "paid", today).
		Select("COALESCE(SUM(paid_amount),0)").Scan(&totalAmount)

	// 今日扫码量（来源为 driver_qr 的订单数）
	var scanCount int64
	database.DB.Model(&models.Order{}).
		Where("source = ? AND DATE(created_at) = ?", "driver_qr", today).
		Count(&scanCount)

	// 扫码转化率
	var fromScanPaid int64
	database.DB.Model(&models.Order{}).
		Where("source = ? AND status = ? AND DATE(created_at) = ?", "driver_qr", "paid", today).
		Count(&fromScanPaid)
	conversionRate := float64(0)
	if scanCount > 0 {
		conversionRate = float64(fromScanPaid) / float64(scanCount) * 100
	}

	// 待核销量
	var pendingVerify int64
	database.DB.Model(&models.Ticket{}).Where("status = ?", "valid").Count(&pendingVerify)

	// 昨日对比
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var yesterdayCount int64
	database.DB.Model(&models.Order{}).
		Where("status = ? AND DATE(created_at) = ?", "paid", yesterday).
		Count(&yesterdayCount)

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"today_orders":      orderCount,
		"today_amount":      totalAmount,
		"scan_count":        scanCount,
		"conversion_rate":   conversionRate,
		"pending_verify":    pendingVerify,
		"yesterday_orders":  yesterdayCount,
	}))
}

// AdminTrend 近 7 日交易趋势
func AdminTrend(c *gin.Context) {
	type DayStat struct {
		Date   string  `json:"date"`
		Orders int64   `json:"orders"`
		Amount float64 `json:"amount"`
	}
	var stats []DayStat
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var cnt int64
		var amt float64
		database.DB.Model(&models.Order{}).
			Where("status = ? AND DATE(created_at) = ?", "paid", d).
			Count(&cnt)
		database.DB.Model(&models.Order{}).
			Where("status = ? AND DATE(created_at) = ?", "paid", d).
			Select("COALESCE(SUM(paid_amount),0)").Scan(&amt)
		stats = append(stats, DayStat{Date: d[5:], Orders: cnt, Amount: amt})
	}
	c.JSON(http.StatusOK, dto.Success(stats))
}

func AdminDriverReview(c *gin.Context) {
	var req struct {
		DriverID uint64 `json:"driver_id" binding:"required"`
		Action   string `json:"action" binding:"required"` // approved / rejected
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if req.Action != "approved" && req.Action != "rejected" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "审核动作不合法"))
		return
	}

	var driver models.Driver
	if err := database.DB.First(&driver, req.DriverID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "司机不存在"))
		return
	}

	tx := database.DB.Begin()
	status := "active"
	if req.Action == "rejected" {
		status = "rejected"
	}
	if err := tx.Model(&driver).Update("status", status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "更新司机状态失败"))
		return
	}
	if err := tx.Model(&models.Vehicle{}).Where("driver_id = ?", driver.ID).Update("status", status).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "更新车辆状态失败"))
		return
	}

	var qrCode string
	if req.Action == "approved" {
		var qr models.DriverQRCode
		err := tx.Where("driver_id = ? AND status = ?", driver.ID, "active").First(&qr).Error
		if err == nil {
			qrCode = qr.Code
		} else {
			qrCode = fmt.Sprintf("DRQR-%s-%d", driver.DriverNo, time.Now().Unix())
			qr = models.DriverQRCode{
				DriverID: driver.ID,
				Code:     qrCode,
				Scene:    "seat",
				Status:   "active",
			}
			if err := tx.Create(&qr).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, dto.Fail(500, "生成司机二维码失败"))
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "提交审核失败"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"driver_id": driver.ID,
		"status":    status,
		"qr_code":   qrCode,
	}))
}

// ==================== 订单管理（后台） ====================

// AdminOrderList 后台订单列表（支持状态筛选、分页）
func AdminOrderList(c *gin.Context) {
	status := c.Query("status")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var orders []models.Order
	q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&orders)

	c.JSON(http.StatusOK, dto.Page(orders, total, page, size))
}

// AdminOrderRefund 后台处理退款
func AdminOrderRefund(c *gin.Context) {
	var req struct {
		OrderID   uint64  `json:"order_id" binding:"required"`
		Reason    string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	tx := database.DB.Begin()

	// 创建退款单
	var refund models.Refund
	now := timeNow()
	database.DB.Model(&models.Order{}).Where("id = ?", req.OrderID).
		Select("paid_amount").Scan(&refund.Amount)

	refund.OrderID = req.OrderID
	refund.RefundNo = "RF" + time.Now().Format("20060102150405")
	if req.Reason != "" {
		refund.Reason = &req.Reason
	}
	refund.Status = "approved"
	refund.ProcessedAt = &now
	tx.Create(&refund)

	// 更新订单状态
	tx.Model(&models.Order{}).Where("id = ?", req.OrderID).Updates(map[string]interface{}{
		"status":       "refunded",
		"cancelled_at": now,
	})

	tx.Commit()
	c.JSON(http.StatusOK, dto.Success(gin.H{"refund_no": refund.RefundNo}))
}

// ==================== 轮播图管理 ====================

// AdminBannerList 轮播图全列表（含草稿/下架）
func AdminBannerList(c *gin.Context) {
	var banners []models.Banner
	database.DB.Order("sort_order").Find(&banners)
	c.JSON(http.StatusOK, dto.Success(banners))
}

// AdminBannerCreate 新增轮播图
func AdminBannerCreate(c *gin.Context) {
	var b models.Banner
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	database.DB.Create(&b)
	c.JSON(http.StatusOK, dto.Success(b))
}

// AdminBannerUpdate 编辑轮播图
func AdminBannerUpdate(c *gin.Context) {
	id := c.Param("id")
	var b models.Banner
	if err := database.DB.First(&b, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "不存在"))
		return
	}
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	database.DB.Save(&b)
	c.JSON(http.StatusOK, dto.Success(b))
}

// AdminBannerDelete 删除轮播图
func AdminBannerDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Banner{}, id)
	c.JSON(http.StatusOK, dto.Success(nil))
}

// ==================== 佣金批次 ====================

// AdminCommissionBatches 佣金发放批次列表
func AdminCommissionBatches(c *gin.Context) {
	type BatchRow struct {
		Status          string  `json:"status"`
		Count           int64   `json:"count"`
		DriverCnt       int64   `json:"driver_count"`
		CommissionTotal float64 `json:"commission_total"`
		SettledDate     string  `json:"settled_date"`
	}
	var rows []BatchRow
	database.DB.Model(&models.DriverCommission{}).
		Select("status, COUNT(*) as count, COUNT(DISTINCT driver_id) as driver_cnt, SUM(commission_amount) as commission_total, DATE(created_at) as settled_date").
		Group("status, settled_date").Order("settled_date DESC").Limit(10).Scan(&rows)
	c.JSON(http.StatusOK, dto.Success(rows))
}

// AdminCommissionSettle 手动触发佣金结算
func AdminCommissionSettle(c *gin.Context) {
	var req struct {
		DriverIDs []uint64 `json:"driver_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	now := timeNow()
	result := database.DB.Model(&models.DriverCommission{}).
		Where("driver_id IN ? AND status = ?", req.DriverIDs, "pending").
		Updates(map[string]interface{}{
			"status":     "settled",
			"settled_at": now,
		})
	c.JSON(http.StatusOK, dto.Success(gin.H{"affected_rows": result.RowsAffected}))
}

// ==================== 参数配置 ====================

// AdminParams 获取系统参数（司机提成比例等）
func AdminParams(c *gin.Context) {
	// 从司机表取默认佣金率，实际项目可建 config 表
	var rate float64
	database.DB.Model(&models.Driver{}).Select("COALESCE(AVG(commission_rate),0.08)").Scan(&rate)
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"commission_rate":       rate,
		"platform_fee_rate":     0.03,
		"refund_freeze_hours":   24,
		"settle_cycle":          "T+1",
		"min_withdrawal":        100.0,
	}))
}

// ==================== 帮助函数 ====================

func queryInt(c *gin.Context, key string, fallback int) int {
	if v := c.Query(key); v != "" {
		n := 0
		for _, r := range v {
			n = n*10 + int(r-'0')
		}
		return n
	}
	return fallback
}
