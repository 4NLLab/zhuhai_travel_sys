package handlers

import (
	"net/http"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"
	"zhuhai_travel_backend/security"

	"github.com/gin-gonic/gin"
)

// DriverList 司机列表
func DriverList(c *gin.Context) {
	var drivers []models.Driver
	database.DB.Find(&drivers)
	for i := range drivers {
		drivers[i].Phone = security.MaskPhone(drivers[i].Phone)
		if drivers[i].IdCardNo != nil {
			masked := security.MaskIDCard(*drivers[i].IdCardNo)
			drivers[i].IdCardNo = &masked
		}
	}
	c.JSON(http.StatusOK, dto.Success(drivers))
}

// DriverCommissionList 司机佣金列表
func DriverCommissionList(c *gin.Context) {
	driverID := c.Query("driver_id")
	var commissions []models.DriverCommission
	q := database.DB.Where("driver_id = ?", driverID).Order("created_at DESC")
	q.Find(&commissions)
	c.JSON(http.StatusOK, dto.Success(commissions))
}

// CommissionSummary 佣金汇总（管理后台用）
func CommissionSummary(c *gin.Context) {
	type Result struct {
		Status    string  `json:"status"`
		Total     float64 `json:"total"`
		Count     int64   `json:"count"`
		DriverCnt int64   `json:"driver_count"`
	}
	var results []Result
	database.DB.Model(&models.DriverCommission{}).
		Select("status, SUM(commission_amount) as total, COUNT(*) as count, COUNT(DISTINCT driver_id) as driver_count").
		Group("status").
		Scan(&results)
	c.JSON(http.StatusOK, dto.Success(results))
}
