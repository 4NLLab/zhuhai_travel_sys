package handlers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

type islandEnvelope struct {
	Header map[string]interface{} `json:"header"`
	Body   map[string]interface{} `json:"body"`
}

var islandClient = &http.Client{Timeout: 12 * time.Second}

type islandRouteCandidate struct {
	UpPortID     int64  `json:"up_port_id"`
	DownPortID   int64  `json:"down_port_id"`
	UpPortName   string `json:"up_port_name"`
	DownPortName string `json:"down_port_name"`
	Label        string `json:"label"`
}

var islandDefaultRoutes = []islandRouteCandidate{
	{UpPortID: 2312, DownPortID: 2332, UpPortName: "湾仔码头", DownPortName: "澳门环岛游A", Label: "澳门环岛游"},
	{UpPortID: 2312, DownPortID: 2504, UpPortName: "湾仔码头", DownPortName: "澳门环岛游B", Label: "澳门环岛游"},
	{UpPortID: 2312, DownPortID: 2317, UpPortName: "湾仔码头", DownPortName: "珠澳夜游A", Label: "珠澳夜游"},
	{UpPortID: 2312, DownPortID: 2505, UpPortName: "湾仔码头", DownPortName: "珠澳夜游B", Label: "珠澳夜游"},
}

func IslandCruisePorts(c *gin.Context) {
	callIslandAPI(c, "/portdistribute/api/basicInfo/searchPort", map[string]interface{}{})
}

func IslandCruiseCertTypes(c *gin.Context) {
	callIslandAPI(c, "/portdistribute/api/basicInfo/searchCertType", map[string]interface{}{})
}

func IslandCruiseVoyages(c *gin.Context) {
	var req struct {
		DepartureDate string `form:"departure_date"`
		UpPortID      int64  `form:"up_port_id"`
		DownPortID    int64  `form:"down_port_id"`
		PeopleNum     int    `form:"people_num"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	if req.DepartureDate == "" {
		req.DepartureDate = time.Now().Format("2006-01-02")
	}
	if req.UpPortID == 0 {
		req.UpPortID = 2312
	}
	if req.DownPortID == 0 {
		req.DownPortID = 2332
	}
	if req.PeopleNum <= 0 {
		req.PeopleNum = 1
	}
	callIslandAPI(c, "/portdistribute/api/voyage/searchVoyage", map[string]interface{}{
		"departureDate": req.DepartureDate,
		"upPortId":      req.UpPortID,
		"downPortId":    req.DownPortID,
		"peopleNum":     req.PeopleNum,
	})
}

func IslandCruiseVoyageCalendar(c *gin.Context) {
	var req struct {
		StartDate  string `form:"start_date"`
		Days       int    `form:"days"`
		UpPortID   int64  `form:"up_port_id"`
		DownPortID int64  `form:"down_port_id"`
		PeopleNum  int    `form:"people_num"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	if req.StartDate == "" {
		req.StartDate = time.Now().Format("2006-01-02")
	}
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "日期格式错误"))
		return
	}
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Days > 15 {
		req.Days = 15
	}
	if req.UpPortID == 0 {
		req.UpPortID = 2312
	}
	if req.DownPortID == 0 {
		req.DownPortID = 2332
	}
	if req.PeopleNum <= 0 {
		req.PeopleNum = 1
	}

	days := make([]gin.H, 0, req.Days)
	for offset := 0; offset < req.Days; offset++ {
		departureDate := start.AddDate(0, 0, offset).Format("2006-01-02")
		parsed, err := islandAPIRequest("/portdistribute/api/voyage/searchVoyage", map[string]interface{}{
			"departureDate": departureDate,
			"upPortId":      req.UpPortID,
			"downPortId":    req.DownPortID,
			"peopleNum":     req.PeopleNum,
		})
		if err != nil {
			days = append(days, gin.H{"date": departureDate, "count": 0, "voyages": []interface{}{}})
			continue
		}
		voyages := interfaceList(parsed["data"])
		days = append(days, gin.H{"date": departureDate, "count": len(voyages), "voyages": voyages})
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{"days": days}))
}

func IslandCruiseSmartSearch(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date"`
		Days      int    `form:"days"`
		PeopleNum int    `form:"people_num"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	if req.StartDate == "" {
		req.StartDate = time.Now().Format("2006-01-02")
	}
	start, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "日期格式错误"))
		return
	}
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Days > 15 {
		req.Days = 15
	}
	if req.PeopleNum <= 0 {
		req.PeopleNum = 1
	}

	type routeResult struct {
		route      islandRouteCandidate
		date       string
		voyages    []interface{}
		minPrice   float64
		firstClock string
	}

	results := make([]routeResult, 0, len(islandDefaultRoutes))
	var best *routeResult
	for _, route := range islandDefaultRoutes {
		routeSummary := routeResult{route: route}
		for offset := 0; offset < req.Days; offset++ {
			departureDate := start.AddDate(0, 0, offset).Format("2006-01-02")
			parsed, err := islandAPIRequest("/portdistribute/api/voyage/searchVoyage", map[string]interface{}{
				"departureDate": departureDate,
				"upPortId":      route.UpPortID,
				"downPortId":    route.DownPortID,
				"peopleNum":     req.PeopleNum,
			})
			if err != nil {
				continue
			}
			voyages := interfaceList(parsed["data"])
			if len(voyages) == 0 {
				continue
			}
			routeSummary.date = departureDate
			routeSummary.voyages = voyages
			routeSummary.minPrice = minVoyagePrice(voyages)
			routeSummary.firstClock = firstVoyageClock(voyages)
			break
		}
		results = append(results, routeSummary)
		if len(routeSummary.voyages) > 0 && best == nil {
			current := routeSummary
			best = &current
		}
	}

	routes := make([]gin.H, 0, len(results))
	for _, result := range results {
		routes = append(routes, gin.H{
			"up_port_id":     result.route.UpPortID,
			"down_port_id":   result.route.DownPortID,
			"up_port_name":   result.route.UpPortName,
			"down_port_name": result.route.DownPortName,
			"label":          result.route.Label,
			"date":           result.date,
			"count":          len(result.voyages),
			"min_price":      result.minPrice,
			"first_time":     result.firstClock,
			"voyages":        result.voyages,
		})
	}
	var recommended gin.H
	if best != nil {
		recommended = gin.H{
			"up_port_id":     best.route.UpPortID,
			"down_port_id":   best.route.DownPortID,
			"up_port_name":   best.route.UpPortName,
			"down_port_name": best.route.DownPortName,
			"label":          best.route.Label,
			"date":           best.date,
			"count":          len(best.voyages),
			"min_price":      best.minPrice,
			"first_time":     best.firstClock,
			"voyages":        best.voyages,
		}
	}
	c.JSON(http.StatusOK, dto.Success(gin.H{"recommended": recommended, "routes": routes}))
}

func IslandCruisePrice(c *gin.Context) {
	var req struct {
		VoyageID   int64 `form:"voyage_id" binding:"required"`
		UpPortID   int64 `form:"up_port_id"`
		DownPortID int64 `form:"down_port_id"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少 voyage_id"))
		return
	}
	if req.UpPortID == 0 {
		req.UpPortID = 2312
	}
	if req.DownPortID == 0 {
		req.DownPortID = 2332
	}
	callIslandAPI(c, "/portdistribute/api/voyage/searchPrice", map[string]interface{}{
		"voyageId":   req.VoyageID,
		"upPortId":   req.UpPortID,
		"downPortId": req.DownPortID,
	})
}

func IslandCruiseBalance(c *gin.Context) {
	parsed, err := islandAPIRequest("/portdistribute/api/preDeposit/searchBalance", map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(islandDataMap(parsed)))
}

type islandLockPassenger struct {
	Name           string  `json:"name"`
	CertTypeID     int64   `json:"certTypeId"`
	CertNo         string  `json:"certNo"`
	Mobile         string  `json:"mobile"`
	CabinClassID   string  `json:"cabinClassId"`
	CabinClassName string  `json:"cabinClassName"`
	CabinTypeCode  string  `json:"cabinTypeCode"`
	FareType       string  `json:"fareType"`
	FareTypeName   string  `json:"fareTypeName"`
	OriginalPrice  float64 `json:"originalPrice"`
	Price          float64 `json:"price"`
	Trip           int     `json:"trip"`
}

func IslandCruiseLockOrder(c *gin.Context) {
	var req struct {
		VoyageID              int64                 `json:"voyageId"`
		UpPortID              int64                 `json:"upPortId"`
		DownPortID            int64                 `json:"downPortId"`
		LinkMan               string                `json:"linkMan"`
		Mobile                string                `json:"mobile"`
		GoTime                string                `json:"goTime"`
		ThirdOrderNo          string                `json:"thirdOrderNo"`
		VoyageName            string                `json:"voyageName"`
		VoyageNo              string                `json:"voyageNo"`
		ShipName              string                `json:"shipName"`
		UpPortName            string                `json:"upPortName"`
		DownPortName          string                `json:"downPortName"`
		PassengerCertInfoList []islandLockPassenger `json:"passengerCertInfoList"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	if req.VoyageID <= 0 || req.UpPortID <= 0 || req.DownPortID <= 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少航班或港口参数"))
		return
	}
	req.LinkMan = strings.TrimSpace(req.LinkMan)
	req.Mobile = strings.TrimSpace(req.Mobile)
	if req.LinkMan == "" || req.Mobile == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少联系人或手机号"))
		return
	}
	if len(req.PassengerCertInfoList) == 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "至少需要 1 位出行人"))
		return
	}
	req.GoTime = strings.TrimSpace(req.GoTime)
	if req.GoTime == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少开航时间"))
		return
	}

	passengers := make([]map[string]interface{}, 0, len(req.PassengerCertInfoList))
	totalAmount := 0.0
	for index, passenger := range req.PassengerCertInfoList {
		normalized, err := normalizeIslandPassenger(passenger)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Fail(400, fmt.Sprintf("第 %d 位出行人%s", index+1, err.Error())))
			return
		}
		totalAmount += passenger.Price
		passengers = append(passengers, normalized)
	}
	localOrderNo := generateIslandOrderNo()
	if strings.TrimSpace(req.ThirdOrderNo) == "" {
		req.ThirdOrderNo = localOrderNo
	}

	body := map[string]interface{}{
		"voyageId":              req.VoyageID,
		"upPortId":              req.UpPortID,
		"downPortId":            req.DownPortID,
		"linkMan":               req.LinkMan,
		"mobile":                req.Mobile,
		"thirdOrderNo":          req.ThirdOrderNo,
		"passengerCertInfoList": passengers,
	}
	if strings.TrimSpace(req.GoTime) != "" {
		body["goTime"] = strings.TrimSpace(req.GoTime)
	}
	lockRequest := mustJSON(body)
	order := models.IslandCruiseOrder{
		LocalOrderNo:   localOrderNo,
		ThirdOrderNo:   req.ThirdOrderNo,
		Status:         "pending_lock",
		VoyageID:       req.VoyageID,
		VoyageName:     strings.TrimSpace(req.VoyageName),
		VoyageNo:       strings.TrimSpace(req.VoyageNo),
		ShipName:       strings.TrimSpace(req.ShipName),
		UpPortID:       req.UpPortID,
		UpPortName:     strings.TrimSpace(req.UpPortName),
		DownPortID:     req.DownPortID,
		DownPortName:   strings.TrimSpace(req.DownPortName),
		GoTime:         req.GoTime,
		ContactName:    req.LinkMan,
		ContactMobile:  req.Mobile,
		PassengerCount: len(req.PassengerCertInfoList),
		TotalAmount:    roundMoney(totalAmount),
		LockRequest:    &lockRequest,
	}
	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建环岛游订单失败"))
		return
	}
	for _, passenger := range req.PassengerCertInfoList {
		record := models.IslandCruisePassenger{
			IslandCruiseOrderID: order.ID,
			Name:                strings.TrimSpace(passenger.Name),
			Mobile:              strings.TrimSpace(passenger.Mobile),
			CertTypeID:          passenger.CertTypeID,
			CertNo:              strings.TrimSpace(passenger.CertNo),
			CabinClassID:        strings.TrimSpace(passenger.CabinClassID),
			CabinClassName:      strings.TrimSpace(passenger.CabinClassName),
			CabinTypeCode:       strings.TrimSpace(passenger.CabinTypeCode),
			FareType:            strings.TrimSpace(passenger.FareType),
			FareTypeName:        strings.TrimSpace(passenger.FareTypeName),
			OriginalPrice:       roundMoney(passenger.OriginalPrice),
			Price:               roundMoney(passenger.Price),
			Trip:                passenger.Trip,
		}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, dto.Fail(500, "创建环岛游乘客失败"))
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "保存环岛游订单失败"))
		return
	}

	parsed, err := islandAPIRequest("/portdistribute/api/order/lockOrderForNet", body)
	if err != nil {
		database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":        "lock_failed",
			"lock_response": mustJSON(map[string]interface{}{"error": err.Error()}),
		})
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}

	data := islandDataMap(parsed)
	supplierOrderNo := stringFromMap(data, "orderNo")
	supplierTicketNo := firstChildString(data, "childOrderInfoList", "ticketNo")
	if supplierOrderNo == "" {
		failPayload := mustJSON(parsed)
		database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":        "lock_failed",
			"lock_response": failPayload,
		})
		c.JSON(http.StatusBadGateway, dto.Fail(502, "供应商未返回有效订单号"))
		return
	}
	now := time.Now()
	expireAt := now.Add(10 * time.Minute)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":             "pending_payment",
		"supplier_order_no":  supplierOrderNo,
		"supplier_ticket_no": supplierTicketNo,
		"lock_response":      mustJSON(parsed),
		"locked_at":          &now,
		"lock_expire_at":     &expireAt,
	})

	c.JSON(http.StatusOK, dto.Success(gin.H{
		"local_order_no": localOrderNo,
		"ticketNo":       supplierTicketNo,
		"payAmount":      order.TotalAmount,
		"lock_expire_at": expireAt.Format(time.RFC3339),
	}))
}

func IslandCruiseSaleOrder(c *gin.Context) {
	var req struct {
		LocalOrderNo  string  `json:"local_order_no"`
		OrderNo       string  `json:"orderNo"`
		PayType       int64   `json:"payType"`
		PayEvidenceNo string  `json:"payEvidenceNo"`
		PayAmount     float64 `json:"payAmount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	req.LocalOrderNo = strings.TrimSpace(req.LocalOrderNo)
	req.OrderNo = strings.TrimSpace(req.OrderNo)
	var order models.IslandCruiseOrder
	query := database.DB.Where("local_order_no = ?", req.LocalOrderNo)
	if req.LocalOrderNo == "" {
		query = database.DB.Where("supplier_order_no = ?", req.OrderNo)
	}
	if err := query.First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return
	}
	if order.Status != "pending_payment" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "当前订单状态不可支付"))
		return
	}
	if order.SupplierOrderNo == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "订单尚未完成座位保留"))
		return
	}
	if order.LockExpireAt != nil && time.Now().After(*order.LockExpireAt) {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "订单已超时，请重新选择班次"))
		return
	}
	if !ensureIslandBalance(c, &order) {
		return
	}
	if req.PayAmount <= 0 {
		req.PayAmount = order.TotalAmount
	}
	if math.Abs(req.PayAmount-order.TotalAmount) > 0.01 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "支付金额与订单金额不一致"))
		return
	}
	if req.PayType == 0 {
		req.PayType = 84
	}
	if req.PayEvidenceNo == "" {
		req.PayEvidenceNo = "WX" + time.Now().Format("20060102150405")
	}
	body := map[string]interface{}{
		"orderNo":       order.SupplierOrderNo,
		"payType":       req.PayType,
		"payEvidenceNo": req.PayEvidenceNo,
		"payAmount":     roundMoney(req.PayAmount),
	}
	parsed, err := islandAPIRequest("/portdistribute/api/order/saleOrderForNet", body)
	if err != nil {
		database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":        "sale_failed",
			"sale_response": mustJSON(map[string]interface{}{"error": err.Error()}),
		})
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	now := time.Now()
	data := islandDataMap(parsed)
	ticketNo, codeContent := syncIslandPassengerTickets(order.ID, data)
	if ticketNo == "" {
		ticketNo = order.SupplierTicketNo
	}
	if codeContent == "" {
		codeContent = order.SupplierCodeContent
	}
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":                "ticketed",
		"pay_amount":            roundMoney(req.PayAmount),
		"pay_type":              req.PayType,
		"pay_evidence_no":       req.PayEvidenceNo,
		"supplier_ticket_no":    ticketNo,
		"supplier_code_content": codeContent,
		"sale_response":         mustJSON(parsed),
		"paid_at":               &now,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"local_order_no": order.LocalOrderNo,
		"ticketNo":       ticketNo,
		"codeContent":    codeContent,
		"paid_at":        now.Format(time.RFC3339),
	}))
}

func IslandCruiseUnlockOrder(c *gin.Context) {
	var req struct {
		LocalOrderNo string `json:"local_order_no"`
		OrderNo      string `json:"orderNo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	req.LocalOrderNo = strings.TrimSpace(req.LocalOrderNo)
	req.OrderNo = strings.TrimSpace(req.OrderNo)
	var order models.IslandCruiseOrder
	query := database.DB.Where("local_order_no = ?", req.LocalOrderNo)
	if req.LocalOrderNo == "" {
		query = database.DB.Where("supplier_order_no = ?", req.OrderNo)
	}
	if err := query.First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return
	}
	if order.Status == "ticketed" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "已出票订单不可取消"))
		return
	}
	if order.SupplierOrderNo != "" {
		if _, err := islandAPIRequest("/portdistribute/api/order/unLockOrderForNet", map[string]interface{}{"orderNo": order.SupplierOrderNo}); err != nil {
			c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
			return
		}
	}
	now := time.Now()
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":       "cancelled",
		"cancelled_at": &now,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{"local_order_no": order.LocalOrderNo, "status": "cancelled"}))
}

func IslandCruiseOrder(c *gin.Context) {
	orderNo := strings.TrimSpace(c.Query("order_no"))
	localOrderNo := strings.TrimSpace(c.Query("local_order_no"))
	if orderNo == "" && localOrderNo == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少订单号"))
		return
	}
	var order models.IslandCruiseOrder
	query := database.DB.Preload("Passengers").Where("local_order_no = ?", localOrderNo)
	if localOrderNo == "" {
		query = database.DB.Preload("Passengers").Where("supplier_order_no = ?", orderNo)
	}
	if err := query.First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return
	}
	if order.SupplierOrderNo != "" {
		if parsed, err := islandAPIRequest("/portdistribute/api/order/searchOrder", map[string]interface{}{"orderNo": order.SupplierOrderNo}); err == nil {
			payload := mustJSON(parsed)
			updates := reconcileIslandOrderUpdates(order, islandDataMap(parsed))
			updates["order_response"] = payload
			database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(updates)
			order.OrderResponse = &payload
			database.DB.Preload("Passengers").Where("id = ?", order.ID).First(&order)
		}
	}
	c.JSON(http.StatusOK, dto.Success(publicIslandOrder(order)))
}

func IslandCruiseRefundFee(c *gin.Context) {
	order, ok := loadIslandOrderForAction(c, true)
	if !ok {
		return
	}
	ticketNos := islandOrderTicketNos(order)
	if len(ticketNos) == 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "订单没有可退票号"))
		return
	}
	body := map[string]interface{}{"orderNo": order.SupplierOrderNo, "ticketNoList": ticketNos}
	parsed, err := islandAPIRequest("/portdistribute/api/order/searchRefundFee", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	data := islandDataMap(parsed)
	refundFee := numberFromMap(data, "totalRefundFee")
	payload := mustJSON(parsed)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"refund_status":       "fee_quoted",
		"refund_fee":          roundMoney(refundFee),
		"refund_amount":       roundMoney(math.Max(0, order.PayAmount-refundFee)),
		"refund_fee_response": payload,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"local_order_no": order.LocalOrderNo,
		"refund_fee":     roundMoney(refundFee),
		"refund_amount":  roundMoney(math.Max(0, order.PayAmount-refundFee)),
		"details":        data["refundFeeDetailDTOList"],
	}))
}

func IslandCruiseRefund(c *gin.Context) {
	var req struct {
		LocalOrderNo string                   `json:"local_order_no"`
		TicketNos    []string                 `json:"ticket_nos"`
		FeeDetails   []map[string]interface{} `json:"refund_fee_detail_list"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	order, ok := findIslandOrder(c, req.LocalOrderNo, "")
	if !ok {
		return
	}
	if order.Status != "ticketed" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "只有已出票订单可申请退票"))
		return
	}
	if order.VerifyStatus == "checked" || order.VerifiedAt != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "已核销订单不可退票"))
		return
	}
	details := sanitizeRefundFeeDetails(req.FeeDetails)
	if len(details) == 0 {
		ticketNos := req.TicketNos
		if len(ticketNos) == 0 {
			ticketNos = islandOrderTicketNos(order)
		}
		feeResp, err := islandAPIRequest("/portdistribute/api/order/searchRefundFee", map[string]interface{}{"orderNo": order.SupplierOrderNo, "ticketNoList": ticketNos})
		if err != nil {
			c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
			return
		}
		data := islandDataMap(feeResp)
		for _, item := range interfaceList(data["refundFeeDetailDTOList"]) {
			if row, ok := item.(map[string]interface{}); ok {
				details = append(details, row)
			}
		}
		details = sanitizeRefundFeeDetails(details)
		payload := mustJSON(feeResp)
		database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Update("refund_fee_response", payload)
	}
	if len(details) == 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少退票手续费明细"))
		return
	}
	refundFlowNo := "RF" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", time.Now().Nanosecond()%10000)
	body := map[string]interface{}{"orderNo": order.SupplierOrderNo, "refundFlowNo": refundFlowNo, "refundFeeDetailList": details}
	parsed, err := islandAPIRequest("/portdistribute/api/order/refundTicket", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	refundFee := refundFeeTotal(details)
	payload := mustJSON(parsed)
	now := time.Now()
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"status":          "refunded",
		"refund_status":   "refunded",
		"refund_flow_no":  refundFlowNo,
		"refund_fee":      roundMoney(refundFee),
		"refund_amount":   roundMoney(math.Max(0, order.PayAmount-refundFee)),
		"refund_response": payload,
		"cancelled_at":    &now,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"local_order_no": order.LocalOrderNo,
		"refund_flow_no": refundFlowNo,
		"refund_fee":     roundMoney(refundFee),
		"refund_amount":  roundMoney(math.Max(0, order.PayAmount-refundFee)),
		"status":         "refunded",
	}))
}

func IslandCruiseChangeFee(c *gin.Context) {
	order, ok := loadIslandOrderForAction(c, true)
	if !ok {
		return
	}
	body := map[string]interface{}{"orderNo": order.SupplierOrderNo, "ticketNoList": islandOrderTicketNos(order)}
	parsed, err := islandAPIRequest("/portdistribute/api/order/searchChangeFee", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	data := islandDataMap(parsed)
	payload := mustJSON(parsed)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"change_status":       "fee_quoted",
		"change_fee":          roundMoney(numberFromMap(data, "changeFee")),
		"change_fee_response": payload,
	})
	c.JSON(http.StatusOK, dto.Success(data))
}

func IslandCruiseChangeVoyages(c *gin.Context) {
	order, ok := loadIslandOrderForAction(c, true)
	if !ok {
		return
	}
	departureDate := strings.TrimSpace(c.Query("departure_date"))
	if departureDate == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少改签日期"))
		return
	}
	body := map[string]interface{}{"departureDate": departureDate, "orderNo": order.SupplierOrderNo, "ticketNoList": islandOrderTicketNos(order)}
	parsed, err := islandAPIRequest("/portdistribute/api/voyage/searchVoyageForChange", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	payload := mustJSON(parsed)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Update("change_voyage_response", payload)
	c.JSON(http.StatusOK, dto.Success(parsed["data"]))
}

func IslandCruiseChangeLock(c *gin.Context) {
	var req struct {
		LocalOrderNo string                   `json:"local_order_no"`
		VoyageID     int64                    `json:"voyageId"`
		Tickets      []map[string]interface{} `json:"tickets"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数错误"))
		return
	}
	order, ok := findIslandOrder(c, req.LocalOrderNo, "")
	if !ok {
		return
	}
	if order.Status != "ticketed" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "只有已出票订单可改签"))
		return
	}
	if req.VoyageID <= 0 || len(req.Tickets) == 0 {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少改签航班或票信息"))
		return
	}
	body := map[string]interface{}{"voyageId": req.VoyageID, "orderNo": order.SupplierOrderNo, "tickets": req.Tickets}
	parsed, err := islandAPIRequest("/portdistribute/api/order/changeOrderLock", body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	data := islandDataMap(parsed)
	changeOrderNo := stringFromMap(data, "orderNo")
	payload := mustJSON(parsed)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"change_status":        "locked",
		"change_order_no":      changeOrderNo,
		"change_price_diff":    roundMoney(numberFromMap(data, "priceDiff")),
		"change_lock_response": payload,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{"local_order_no": order.LocalOrderNo, "change_order_no": changeOrderNo, "price_diff": numberFromMap(data, "priceDiff"), "tickets": data["tickets"]}))
}

func IslandCruiseChangeUnlock(c *gin.Context) {
	order, ok := loadIslandOrderForAction(c, false)
	if !ok {
		return
	}
	if order.ChangeOrderNo == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "没有已锁定的改签订单"))
		return
	}
	parsed, err := islandAPIRequest("/portdistribute/api/order/changeOrderUnLock", map[string]interface{}{"orderNo": order.ChangeOrderNo})
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	payload := mustJSON(parsed)
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"change_status":          "unlocked",
		"change_unlock_response": payload,
	})
	c.JSON(http.StatusOK, dto.Success(gin.H{"local_order_no": order.LocalOrderNo, "status": "unlocked"}))
}

func IslandCruiseVerifyNotify(c *gin.Context) {
	var req struct {
		OrderNo string `json:"orderNo"`
		Tickets []struct {
			TicketNo        string `json:"ticketNo"`
			CheckStatus     int    `json:"checkStatus"`
			CheckStatusDesc string `json:"checkStatusDesc"`
		} `json:"tickets"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"IsSuccess": false, "Message": "参数错误"})
		return
	}
	if strings.TrimSpace(req.OrderNo) == "" || len(req.Tickets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"IsSuccess": false, "Message": "缺少订单或票信息"})
		return
	}
	var order models.IslandCruiseOrder
	if err := database.DB.Where("supplier_order_no = ? OR local_order_no = ?", req.OrderNo, req.OrderNo).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"IsSuccess": false, "Message": "订单不存在"})
		return
	}
	allChecked := true
	for _, ticket := range req.Tickets {
		if ticket.CheckStatus != 1 {
			allChecked = false
			break
		}
	}
	status := "partial"
	if allChecked {
		status = "checked"
	}
	payload := mustJSON(req)
	now := time.Now()
	updates := map[string]interface{}{"verify_status": status, "verify_response": payload}
	if allChecked {
		updates["status"] = "used"
		updates["verified_at"] = &now
	}
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"IsSuccess": true, "Message": "核销通知已记录"})
}

func normalizeIslandPassenger(passenger islandLockPassenger) (map[string]interface{}, error) {
	passenger.Name = strings.TrimSpace(passenger.Name)
	passenger.CertNo = strings.TrimSpace(passenger.CertNo)
	passenger.Mobile = strings.TrimSpace(passenger.Mobile)
	passenger.CabinClassID = strings.TrimSpace(passenger.CabinClassID)
	passenger.CabinTypeCode = strings.TrimSpace(passenger.CabinTypeCode)
	if passenger.CabinTypeCode == "" || passenger.CabinTypeCode == "K" {
		passenger.CabinTypeCode = "S"
	}
	passenger.FareType = strings.TrimSpace(passenger.FareType)
	if passenger.Name == "" || passenger.CertNo == "" || passenger.Mobile == "" {
		return nil, fmt.Errorf("实名信息不完整")
	}
	if passenger.CertTypeID <= 0 {
		return nil, fmt.Errorf("缺少证件类型")
	}
	if passenger.CabinClassID == "" || passenger.CabinTypeCode == "" || passenger.FareType == "" {
		return nil, fmt.Errorf("缺少舱位票价参数")
	}
	if passenger.Price < 0 {
		return nil, fmt.Errorf("票价不能小于 0")
	}
	return map[string]interface{}{
		"name":          passenger.Name,
		"certTypeId":    passenger.CertTypeID,
		"certNo":        passenger.CertNo,
		"mobile":        passenger.Mobile,
		"cabinClassId":  passenger.CabinClassID,
		"cabinTypeCode": passenger.CabinTypeCode,
		"fareType":      passenger.FareType,
		"originalPrice": roundMoney(passenger.OriginalPrice),
		"price":         roundMoney(passenger.Price),
		"trip":          passenger.Trip,
	}, nil
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func callIslandAPI(c *gin.Context, path string, body map[string]interface{}) {
	parsed, err := islandAPIRequest(path, body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(parsed))
}

func islandAPIRequest(path string, body map[string]interface{}) (map[string]interface{}, error) {
	cfg := config.Load()
	if cfg.IslandDistributor == "" || cfg.IslandAccessToken == "" {
		return nil, fmt.Errorf("环岛游接口账号未配置")
	}

	header := map[string]interface{}{
		"distributorCode":      cfg.IslandDistributor,
		"ticketMachineAccount": "",
		"timestamp":            time.Now().UnixMilli(),
		"terminalNo":           "",
		"locatePort":           "",
	}
	payload := islandEnvelope{Header: header, Body: body}
	payload.Header["secretKey"] = islandSecret(payload, cfg.IslandAccessToken)

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("环岛游请求序列化失败")
	}

	resp, err := islandClient.Post(cfg.IslandCruiseBaseURL+path, "application/json", bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("环岛游接口请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("环岛游接口响应读取失败")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("环岛游接口 HTTP %d", resp.StatusCode)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("环岛游接口返回非 JSON")
	}
	code := strings.TrimSpace(fmt.Sprint(parsed["code"]))
	if code != "" && code != "200" {
		msg := strings.TrimSpace(fmt.Sprint(parsed["msg"]))
		if msg == "" {
			msg = "环岛游接口业务失败"
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return parsed, nil
}

func generateIslandOrderNo() string {
	now := time.Now()
	return "HD" + now.Format("20060102150405") + fmt.Sprintf("%09d", now.Nanosecond())
}

func mustJSON(value interface{}) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func islandDataMap(parsed map[string]interface{}) map[string]interface{} {
	data, _ := parsed["data"].(map[string]interface{})
	return data
}

func stringFromMap(data map[string]interface{}, key string) string {
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func firstChildString(data map[string]interface{}, listKey string, valueKey string) string {
	list, _ := data[listKey].([]interface{})
	if len(list) == 0 {
		return ""
	}
	first, _ := list[0].(map[string]interface{})
	return stringFromMap(first, valueKey)
}

func syncIslandPassengerTickets(orderID uint64, data map[string]interface{}) (string, string) {
	ticketItems := ticketInfoItems(data)
	if len(ticketItems) == 0 {
		return "", ""
	}
	var passengers []models.IslandCruisePassenger
	if err := database.DB.Where("island_cruise_order_id = ?", orderID).Order("id ASC").Find(&passengers).Error; err != nil {
		return firstTicketNo(ticketItems), firstCodeContent(ticketItems)
	}
	used := make(map[int]bool)
	for _, passenger := range passengers {
		ticketIndex := findTicketForPassenger(ticketItems, used, passenger)
		if ticketIndex < 0 {
			continue
		}
		used[ticketIndex] = true
		ticketNo := stringFromMap(ticketItems[ticketIndex], "ticketNo")
		if ticketNo == "" {
			continue
		}
		codeContent := stringFromMap(ticketItems[ticketIndex], "codeContent")
		if codeContent == "" {
			codeContent = ticketNo
		}
		database.DB.Model(&models.IslandCruisePassenger{}).Where("id = ?", passenger.ID).Updates(map[string]interface{}{
			"supplier_ticket_no":    ticketNo,
			"supplier_code_content": codeContent,
		})
	}
	return firstTicketNo(ticketItems), firstCodeContent(ticketItems)
}

func ticketInfoItems(data map[string]interface{}) []map[string]interface{} {
	keys := []string{"ticketInfoList", "childOrderInfoList", "tickets", "backTicketInfoList"}
	items := make([]map[string]interface{}, 0)
	for _, key := range keys {
		for _, raw := range interfaceList(data[key]) {
			if item, ok := raw.(map[string]interface{}); ok {
				items = append(items, item)
			}
		}
	}
	return items
}

func firstTicketNo(items []map[string]interface{}) string {
	for _, item := range items {
		if ticketNo := stringFromMap(item, "ticketNo"); ticketNo != "" {
			return ticketNo
		}
	}
	return ""
}

func firstCodeContent(items []map[string]interface{}) string {
	for _, item := range items {
		if codeContent := stringFromMap(item, "codeContent"); codeContent != "" {
			return codeContent
		}
	}
	return firstTicketNo(items)
}

func findTicketForPassenger(items []map[string]interface{}, used map[int]bool, passenger models.IslandCruisePassenger) int {
	for index, item := range items {
		if used[index] {
			continue
		}
		if stringFromMap(item, "certNo") == passenger.CertNo || stringFromMap(item, "certificateNo") == passenger.CertNo {
			return index
		}
		if stringFromMap(item, "passengerName") == passenger.Name || stringFromMap(item, "name") == passenger.Name {
			return index
		}
	}
	for index := range items {
		if !used[index] {
			return index
		}
	}
	return -1
}

func reconcileIslandOrderUpdates(order models.IslandCruiseOrder, data map[string]interface{}) map[string]interface{} {
	updates := map[string]interface{}{}
	if ticketNo, codeContent := syncIslandPassengerTickets(order.ID, data); ticketNo != "" {
		updates["supplier_ticket_no"] = ticketNo
		if codeContent != "" {
			updates["supplier_code_content"] = codeContent
		}
	}
	statusCode := int(numberFromMap(data, "orderStatus"))
	switch statusCode {
	case 3:
		updates["status"] = "cancelled"
	case 5:
		updates["status"] = "ticketed"
	case 6:
		updates["status"] = "used"
	case 8:
		updates["status"] = "pending_payment"
	case 9:
		updates["change_status"] = "changed"
	}

	tickets := ticketInfoItems(data)
	if len(tickets) > 0 {
		allChecked := true
		allRefunded := true
		for _, ticket := range tickets {
			if int(numberFromMap(ticket, "checkStatus")) != 1 && int(numberFromMap(ticket, "checkState")) != 1 {
				allChecked = false
			}
			if int(numberFromMap(ticket, "refundState")) != 1 {
				allRefunded = false
			}
		}
		now := time.Now()
		if allChecked {
			updates["verify_status"] = "checked"
			updates["verified_at"] = &now
			updates["status"] = "used"
		} else {
			updates["verify_status"] = "pending"
		}
		if allRefunded {
			updates["refund_status"] = "refunded"
			updates["status"] = "refunded"
			updates["cancelled_at"] = &now
		}
	}
	return updates
}

func ensureIslandBalance(c *gin.Context, order *models.IslandCruiseOrder) bool {
	parsed, err := islandAPIRequest("/portdistribute/api/preDeposit/searchBalance", map[string]interface{}{})
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, "余额检查失败: "+err.Error()))
		return false
	}
	data := islandDataMap(parsed)
	balance := numberFromMap(data, "balance")
	payload := mustJSON(parsed)
	now := time.Now()
	database.DB.Model(&models.IslandCruiseOrder{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
		"balance_amount":     balance,
		"balance_checked_at": &now,
		"balance_response":   payload,
	})
	if balance < order.TotalAmount {
		c.JSON(http.StatusBadRequest, dto.Fail(400, fmt.Sprintf("账户余额不足，可用余额 %.2f", balance)))
		return false
	}
	return true
}

func loadIslandOrderForAction(c *gin.Context, requireTicketed bool) (models.IslandCruiseOrder, bool) {
	localOrderNo := strings.TrimSpace(c.Query("local_order_no"))
	orderNo := strings.TrimSpace(c.Query("order_no"))
	return findIslandOrder(c, localOrderNo, orderNo, requireTicketed)
}

func findIslandOrder(c *gin.Context, localOrderNo string, orderNo string, requireTicketed ...bool) (models.IslandCruiseOrder, bool) {
	var order models.IslandCruiseOrder
	localOrderNo = strings.TrimSpace(localOrderNo)
	orderNo = strings.TrimSpace(orderNo)
	if localOrderNo == "" && orderNo == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少订单号"))
		return order, false
	}
	query := database.DB.Preload("Passengers").Where("local_order_no = ?", localOrderNo)
	if localOrderNo == "" {
		query = database.DB.Preload("Passengers").Where("supplier_order_no = ?", orderNo)
	}
	if err := query.First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "订单不存在"))
		return order, false
	}
	if len(requireTicketed) > 0 && requireTicketed[0] && order.Status != "ticketed" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "订单尚未出票"))
		return order, false
	}
	return order, true
}

func islandOrderTicketNos(order models.IslandCruiseOrder) []string {
	seen := map[string]bool{}
	var result []string
	if order.SupplierTicketNo != "" {
		seen[order.SupplierTicketNo] = true
		result = append(result, order.SupplierTicketNo)
	}
	for _, passenger := range order.Passengers {
		if passenger.SupplierTicketNo != "" && !seen[passenger.SupplierTicketNo] {
			seen[passenger.SupplierTicketNo] = true
			result = append(result, passenger.SupplierTicketNo)
		}
	}
	return result
}

func interfaceList(value interface{}) []interface{} {
	list, _ := value.([]interface{})
	return list
}

func minVoyagePrice(voyages []interface{}) float64 {
	minPrice := 0.0
	for _, rawVoyage := range voyages {
		voyage, ok := rawVoyage.(map[string]interface{})
		if !ok {
			continue
		}
		for _, rawCabin := range interfaceList(voyage["cabinPriceList"]) {
			cabin, ok := rawCabin.(map[string]interface{})
			if !ok {
				continue
			}
			for _, rawFare := range interfaceList(cabin["fareTypeList"]) {
				fare, ok := rawFare.(map[string]interface{})
				if !ok {
					continue
				}
				price := numberFromMap(fare, "price")
				if price <= 0 {
					continue
				}
				if minPrice == 0 || price < minPrice {
					minPrice = price
				}
			}
		}
	}
	return roundMoney(minPrice)
}

func firstVoyageClock(voyages []interface{}) string {
	for _, rawVoyage := range voyages {
		voyage, ok := rawVoyage.(map[string]interface{})
		if !ok {
			continue
		}
		if value := stringFromMap(voyage, "departureTime"); value != "" {
			return value
		}
		if value := stringFromMap(voyage, "upTime"); value != "" {
			return value
		}
	}
	return ""
}

func numberFromMap(data map[string]interface{}, key string) float64 {
	value, ok := data[key]
	if !ok || value == nil {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		number, _ := typed.Float64()
		return number
	default:
		var number float64
		fmt.Sscanf(fmt.Sprint(typed), "%f", &number)
		return number
	}
}

func refundFeeTotal(details []map[string]interface{}) float64 {
	total := 0.0
	for _, item := range details {
		total += numberFromMap(item, "refundFee")
	}
	return total
}

func sanitizeRefundFeeDetails(details []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(details))
	for _, item := range details {
		ticketNo := stringFromMap(item, "ticketNo")
		if ticketNo == "" {
			continue
		}
		row := map[string]interface{}{
			"ticketNo":  ticketNo,
			"refundFee": roundMoney(numberFromMap(item, "refundFee")),
		}
		if orderDetailID := numberFromMap(item, "orderDetailId"); orderDetailID > 0 {
			row["orderDetailId"] = orderDetailID
		}
		result = append(result, row)
	}
	return result
}

func publicIslandOrder(order models.IslandCruiseOrder) gin.H {
	passengers := make([]gin.H, 0, len(order.Passengers))
	for _, passenger := range order.Passengers {
		passengers = append(passengers, gin.H{
			"name":           passenger.Name,
			"mobile":         passenger.Mobile,
			"cert_type_id":   passenger.CertTypeID,
			"cert_no":        maskCertNo(passenger.CertNo),
			"fare_type_name": passenger.FareTypeName,
			"cabin_class":    passenger.CabinClassName,
			"price":          passenger.Price,
			"ticket_no":      passenger.SupplierTicketNo,
			"code_content":   passenger.SupplierCodeContent,
		})
	}
	return gin.H{
		"local_order_no":  order.LocalOrderNo,
		"status":          order.Status,
		"voyage_name":     order.VoyageName,
		"voyage_no":       order.VoyageNo,
		"ship_name":       order.ShipName,
		"up_port_name":    order.UpPortName,
		"down_port_name":  order.DownPortName,
		"go_time":         order.GoTime,
		"contact_name":    order.ContactName,
		"contact_mobile":  order.ContactMobile,
		"passenger_count": order.PassengerCount,
		"total_amount":    order.TotalAmount,
		"pay_amount":      order.PayAmount,
		"refund_status":   order.RefundStatus,
		"refund_fee":      order.RefundFee,
		"refund_amount":   order.RefundAmount,
		"change_status":   order.ChangeStatus,
		"change_order_no": order.ChangeOrderNo,
		"verify_status":   order.VerifyStatus,
		"ticket_no":       order.SupplierTicketNo,
		"code_content":    order.SupplierCodeContent,
		"lock_expire_at":  order.LockExpireAt,
		"paid_at":         order.PaidAt,
		"verified_at":     order.VerifiedAt,
		"passengers":      passengers,
	}
}

func maskCertNo(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 8 {
		return value
	}
	return value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
}

func islandSecret(payload islandEnvelope, accessToken string) string {
	signPayload := islandEnvelope{
		Header: map[string]interface{}{},
		Body:   payload.Body,
	}
	for key, value := range payload.Header {
		if key == "secretKey" {
			continue
		}
		signPayload.Header[key] = value
	}
	signString := canonicalJSON(signPayload) + accessToken
	sum := md5.Sum([]byte(signString))
	return hex.EncodeToString(sum[:])
}

func canonicalJSON(value interface{}) string {
	switch typed := value.(type) {
	case islandEnvelope:
		return `{"body":` + canonicalJSON(typed.Body) + `,"header":` + canonicalJSON(typed.Header) + `}`
	case map[string]interface{}:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var b strings.Builder
		b.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			keyJSON, _ := json.Marshal(key)
			b.Write(keyJSON)
			b.WriteByte(':')
			b.WriteString(canonicalJSON(typed[key]))
		}
		b.WriteByte('}')
		return b.String()
	case []interface{}:
		var b strings.Builder
		b.WriteByte('[')
		for i, item := range typed {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(canonicalJSON(item))
		}
		b.WriteByte(']')
		return b.String()
	default:
		raw, _ := json.Marshal(typed)
		return string(raw)
	}
}
