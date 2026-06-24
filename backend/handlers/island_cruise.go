package handlers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/dto"

	"github.com/gin-gonic/gin"
)

type islandEnvelope struct {
	Header map[string]interface{} `json:"header"`
	Body   map[string]interface{} `json:"body"`
}

var islandClient = &http.Client{Timeout: 12 * time.Second}

func IslandCruisePorts(c *gin.Context) {
	callIslandAPI(c, "/portdistribute/api/basicInfo/searchPort", map[string]interface{}{})
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

func callIslandAPI(c *gin.Context, path string, body map[string]interface{}) {
	cfg := config.Load()
	if cfg.IslandDistributor == "" || cfg.IslandAccessToken == "" {
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "环岛游接口账号未配置"))
		return
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
		c.JSON(http.StatusInternalServerError, dto.Fail(500, "环岛游请求序列化失败"))
		return
	}

	resp, err := islandClient.Post(cfg.IslandCruiseBaseURL+path, "application/json", bytes.NewReader(raw))
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, "环岛游接口请求失败: "+err.Error()))
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, "环岛游接口响应读取失败"))
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusBadGateway, dto.Fail(502, fmt.Sprintf("环岛游接口 HTTP %d", resp.StatusCode)))
		return
	}

	var parsed interface{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		c.JSON(http.StatusBadGateway, dto.Fail(502, "环岛游接口返回非 JSON"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(parsed))
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
