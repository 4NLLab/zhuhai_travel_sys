package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"
	"zhuhai_travel_backend/security"

	"github.com/gin-gonic/gin"
)

// ==================== 司机端 ====================

func DriverRegister(c *gin.Context) {
	var req struct {
		Name              string `json:"name" binding:"required"`
		Phone             string `json:"phone" binding:"required"`
		Password          string `json:"password" binding:"required,min=8"`
		IDCardNo          string `json:"id_card_no" binding:"required"`
		PlateNo           string `json:"plate_no" binding:"required"`
		Model             string `json:"model" binding:"required"`
		Seats             int    `json:"seats" binding:"required,min=1"`
		ServiceCity       string `json:"service_city"`
		IdCardFrontURL    string `json:"id_card_front_url" binding:"required"`
		DriverLicenseURL  string `json:"driver_license_url" binding:"required"`
		VehicleLicenseURL string `json:"vehicle_license_url" binding:"required"`
		VehiclePhotoURL   string `json:"vehicle_photo_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "密码至少 8 位"))
		return
	}

	var exists int64
	database.DB.Model(&models.Driver{}).Where("phone = ?", req.Phone).Count(&exists)
	if exists > 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "该手机号已注册司机"))
		return
	}
	database.DB.Model(&models.Vehicle{}).Where("plate_no = ?", req.PlateNo).Count(&exists)
	if exists > 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "该车牌号已绑定"))
		return
	}

	tx := database.DB.Begin()
	driverNo := fmt.Sprintf("DR-ZH-%s%04d", time.Now().Format("0102"), rand.Intn(10000))
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建密码失败"))
		return
	}
	driver := models.Driver{
		DriverNo:          driverNo,
		Name:              req.Name,
		Phone:             req.Phone,
		PasswordHash:      passwordHash,
		IdCardNo:          &req.IDCardNo,
		Status:            "pending_review",
		IdCardFrontURL:    req.IdCardFrontURL,
		DriverLicenseURL:  req.DriverLicenseURL,
		VehicleLicenseURL: req.VehicleLicenseURL,
		VehiclePhotoURL:   req.VehiclePhotoURL,
		CommissionRate:    0.08,
	}
	if err := tx.Create(&driver).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建司机失败"))
		return
	}
	vehicle := models.Vehicle{
		DriverID: driver.ID,
		PlateNo:  req.PlateNo,
		Model:    &req.Model,
		Seats:    &req.Seats,
		Status:   "pending_review",
	}
	if err := tx.Create(&vehicle).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建车辆失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "提交注册失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"driver_id":  driver.ID,
		"driver_no":  driverNo,
		"vehicle_id": vehicle.ID,
		"status":     "pending_review",
		"message":    "注册申请已提交，等待后台审核",
	}))
}

// DriverUploadLicense 上传司机注册审核图片
func DriverUploadLicense(c *gin.Context) {
	kind := c.PostForm("kind")
	allowedKinds := map[string]bool{
		"id_front":        true,
		"driver_license":  true,
		"vehicle_license": true,
		"vehicle_photo":   true,
	}
	if !allowedKinds[kind] {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "证照类型不合法"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请选择图片"))
		return
	}
	if file.Size > 8*1024*1024 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "图片不能超过 8MB"))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".heic": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "仅支持 JPG / PNG / WEBP / HEIC"))
		return
	}

	day := time.Now().Format("20060102")
	dir := filepath.Join("uploads", "drivers", day)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建上传目录失败"))
		return
	}

	filename := fmt.Sprintf("%s-%d-%04d%s", kind, time.Now().UnixNano(), rand.Intn(10000), ext)
	dst := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "保存图片失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"url": "/" + filepath.ToSlash(dst),
	}))
}

// DriverLogin 司机登录（手机号 + 密码）
func DriverLogin(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	var driver models.Driver
	if err := database.DB.Where("phone = ? AND status = ?", req.Phone, "active").First(&driver).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "司机不存在、未审核通过或已禁用"))
		return
	}
	if driver.PasswordHash == "" || !security.VerifyPassword(req.Password, driver.PasswordHash) {
		c.JSON(http.StatusUnauthorized, dto.Fail(401, "手机号或密码错误"))
		return
	}

	token, err := security.GenerateToken(config.Load().JWTSecret, driver.ID, "driver", driver.Name, tokenTTL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "生成登录凭证失败"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"driver_id":       driver.ID,
		"driver_no":       driver.DriverNo,
		"name":            driver.Name,
		"commission_rate": driver.CommissionRate,
		"access_token":    token,
		"token_type":      "Bearer",
		"expires_in":      int(tokenTTL().Seconds()),
	}))
}

// DriverWallet 司机钱包余额
func DriverWallet(c *gin.Context) {
	driverID := currentUserID(c)

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

	// 已提现总额（已打款）
	var withdrawnAmt float64
	database.DB.Model(&models.DriverWithdrawal{}).
		Where("driver_id = ? AND status IN ?", driverID, []string{"approved", "transferred"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&withdrawnAmt)

	// 审核中提现
	var frozenAmt float64
	database.DB.Model(&models.DriverWithdrawal{}).
		Where("driver_id = ? AND status = ?", driverID, "pending").
		Select("COALESCE(SUM(amount), 0)").Scan(&frozenAmt)

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"pending_total":   pendingAmt,
		"settled_total":   settledAmt,
		"withdrawn_total": withdrawnAmt,
		"frozen_total":    frozenAmt,
		"available":       settledAmt - withdrawnAmt - frozenAmt,
	}))
}

// DriverCommissionList 司机佣金明细（可筛选状态）
func DriverViewCommissionList(c *gin.Context) {
	driverID := currentUserID(c)
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
		Amount   float64 `json:"amount" binding:"required,min=0.01"`
		Channel  string  `json:"channel"`   // alipay / wechat
		Account  string  `json:"account"`   // 支付宝账号
		RealName string  `json:"real_name"` // 实名
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "参数错误"))
		return
	}

	if req.Channel == "" {
		req.Channel = "alipay"
	}
	if req.Channel != "alipay" && req.Channel != "wechat" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "提现渠道不支持"))
		return
	}
	driverID := currentUserID(c)

	tx := database.DB.Begin()
	// 计算可提现余额 — pending 状态的提现也要扣除（防并发超提）
	var settledAmt float64
	tx.Model(&models.DriverCommission{}).
		Where("driver_id = ? AND status = ?", driverID, "settled").
		Select("COALESCE(SUM(commission_amount), 0)").Scan(&settledAmt)

	var withdrawnAmt float64
	tx.Model(&models.DriverWithdrawal{}).
		Where("driver_id = ? AND status IN ?", driverID,
			[]string{"pending", "approved", "transferred"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&withdrawnAmt)

	available := settledAmt - withdrawnAmt
	if req.Amount > available {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, dto.Fail(400,
			fmt.Sprintf("余额不足，可提现: ¥%.2f", available)))
		return
	}

	// 创建提现记录
	wdNo := fmt.Sprintf("WD%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
	wd := models.DriverWithdrawal{
		DriverID:     driverID,
		WithdrawalNo: wdNo,
		Amount:       req.Amount,
		Channel:      req.Channel,
		Account:      &req.Account,
		RealName:     &req.RealName,
	}
	if err := tx.Create(&wd).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "提交失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
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
	driverID := currentUserID(c)
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
	if req.Action != "approved" && req.Action != "transferred" && req.Action != "rejected" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "提现操作不支持"))
		return
	}

	var wd models.DriverWithdrawal
	if err := database.DB.First(&wd, req.WithdrawalID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "提现记录不存在"))
		return
	}
	if (req.Action == "approved" || req.Action == "rejected") && wd.Status != "pending" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "只有待审核提现可审核"))
		return
	}
	if req.Action == "transferred" && wd.Status != "approved" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "只有已通过提现可标记打款"))
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
