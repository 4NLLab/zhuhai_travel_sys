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

// ==================== 司机端 ====================

// DriverLogin 司机登录（手机号验证，生产环境需加验证码/密码）
func DriverLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var driver models.Driver
	if err := database.DB.Where("phone = ? AND status = ?", req.Phone, "active").First(&driver).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "司机不存在或已禁用"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"driver_id":       driver.ID,
		"driver_no":       driver.DriverNo,
		"name":            driver.Name,
		"commission_rate": driver.CommissionRate,
	}))
}

// DriverWallet 司机钱包余额
func DriverWallet(c *gin.Context) {
	driverID := c.Query("driver_id")

	// 待结算佣金
	var pendingAmt float64
	database.DB.Model(&models.DriverCommission{}).
		Where("driver_id = ? AND status = ?", driverID, "pending").
		Select("COALESCE(SUM(commission_amount), 0)").Scan(&pendingAmt)

	// 已结算佣金
	var settledAmt float64
	database.DB.Model(&models.DriverCommission{}).
		Where("driver_id = ? AND status = ?", driverID, "settled").
		Select("COALESCE(SUM(commission_amount), 0)").Scan(&settledAmt)

	// 已提现总额
	var withdrawnAmt float64
	database.DB.Model(&models.DriverWithdrawal{}).
		Where("driver_id = ? AND status IN ?", driverID, []string{"approved", "transferred"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&withdrawnAmt)

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"pending_total":   pendingAmt,
		"settled_total":   settledAmt,
		"withdrawn_total": withdrawnAmt,
		"available":       settledAmt - withdrawnAmt, // 可提现余额
	}))
}

// DriverCommissionList 司机佣金明细（可筛选状态）
func DriverViewCommissionList(c *gin.Context) {
	driverID := c.Query("driver_id")
	status := c.Query("status")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.DriverCommission{}).Where("driver_id = ?", driverID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var items []models.DriverCommission
	q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items)

	// 构建返回数据，附加订单号方便司机核对
	type CommissionItem struct {
		models.DriverCommission
		OrderNo string `json:"order_no"`
	}
	var result []CommissionItem
	for _, c := range items {
		var order models.Order
		database.DB.Select("order_no").First(&order, c.OrderID)
		result = append(result, CommissionItem{DriverCommission: c, OrderNo: order.OrderNo})
	}

	c.JSON(http.StatusOK, dto.Page(result, total, page, size))
}

// DriverWithdraw 司机发起提现
func DriverWithdraw(c *gin.Context) {
	var req struct {
		DriverID uint64  `json:"driver_id" binding:"required"`
		Amount   float64 `json:"amount" binding:"required,min=0.01"`
		Channel  string  `json:"channel"`  // alipay / wechat
		Account  string  `json:"account"`  // 支付宝账号
		RealName string  `json:"real_name"` // 实名
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	if req.Channel == "" {
		req.Channel = "alipay"
	}

	// 计算可提现余额
	var settledAmt float64
	database.DB.Model(&models.DriverCommission{}).
		Where("driver_id = ? AND status = ?", req.DriverID, "settled").
		Select("COALESCE(SUM(commission_amount), 0)").Scan(&settledAmt)

	var withdrawnAmt float64
	database.DB.Model(&models.DriverWithdrawal{}).
		Where("driver_id = ? AND status IN ?", req.DriverID, []string{"approved", "transferred"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&withdrawnAmt)

	available := settledAmt - withdrawnAmt
	if req.Amount > available {
		c.JSON(http.StatusBadRequest, dto.Fail(400,
			fmt.Sprintf("余额不足，可提现: ¥%.2f", available)))
		return
	}

	// 创建提现记录
	wdNo := fmt.Sprintf("WD%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
	wd := models.DriverWithdrawal{
		DriverID:     req.DriverID,
		WithdrawalNo: wdNo,
		Amount:       req.Amount,
		Channel:      req.Channel,
		Account:      &req.Account,
		RealName:     &req.RealName,
	}
	if err := database.DB.Create(&wd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "提交失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"withdrawal_no": wdNo,
		"amount":        req.Amount,
		"status":        "pending",
		"message":       "提现申请已提交，待管理员审核打款",
	}))
}

// DriverWithdrawalHistory 司机提现记录
func DriverWithdrawalHistory(c *gin.Context) {
	driverID := c.Query("driver_id")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.DriverWithdrawal{}).Where("driver_id = ?", driverID)

	var total int64
	q.Count(&total)

	var items []models.DriverWithdrawal
	q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items)

	c.JSON(http.StatusOK, dto.Page(items, total, page, size))
}

// ==================== 管理员处理提现 ====================

// AdminWithdrawalList 后台提现列表（待审核）
func AdminWithdrawalList(c *gin.Context) {
	status := c.Query("status")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.DriverWithdrawal{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var items []models.DriverWithdrawal
	q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items)

	c.JSON(http.StatusOK, dto.Page(items, total, page, size))
}

// AdminWithdrawalProcess 管理员处理提现（通过/拒绝/标记已转账）
func AdminWithdrawalProcess(c *gin.Context) {
	var req struct {
		WithdrawalID uint64 `json:"withdrawal_id" binding:"required"`
		Action       string `json:"action" binding:"required"` // approved / transferred / rejected
		AdminID      uint64 `json:"admin_id"`
		Remark       string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var wd models.DriverWithdrawal
	if err := database.DB.First(&wd, req.WithdrawalID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "提现记录不存在"))
		return
	}
	if wd.Status != "pending" && req.Action == "approved" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "该记录已处理"))
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       req.Action,
		"processed_at": now,
	}
	if req.AdminID > 0 {
		updates["processed_by"] = req.AdminID
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	database.DB.Model(&wd).Updates(updates)

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"withdrawal_no": wd.WithdrawalNo,
		"status":        req.Action,
	}))
}
