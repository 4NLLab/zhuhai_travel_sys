package handlers

import (
	"encoding/json"
	"net/http"

	"zhuhai_travel_backend/config"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/services/jzg"

	"github.com/gin-gonic/gin"
)

type cruiseEndpoint struct {
	Path       string
	AuthRole   string
	AllowEmpty bool
}

var cruiseEndpoints = map[string]cruiseEndpoint{
	"ports":          {Path: "/portdistribute/api/basicInfo/searchPort", AllowEmpty: true},
	"cert-types":     {Path: "/portdistribute/api/basicInfo/searchCertType", AllowEmpty: true},
	"voyages":        {Path: "/portdistribute/api/voyage/searchVoyage"},
	"prices":         {Path: "/portdistribute/api/voyage/searchPrice"},
	"single-voyage":  {Path: "/portdistribute/api/voyage/searchSingleVoyage"},
	"lock-order":     {Path: "/portdistribute/api/order/lockOrderForNet", AuthRole: "user"},
	"sale-order":     {Path: "/portdistribute/api/order/saleOrderForNet", AuthRole: "user"},
	"unlock-order":   {Path: "/portdistribute/api/order/unLockOrderForNet", AuthRole: "user"},
	"order":          {Path: "/portdistribute/api/order/searchOrder", AuthRole: "user"},
	"refund-fee":     {Path: "/portdistribute/api/order/searchRefundFee", AuthRole: "user"},
	"refund-ticket":  {Path: "/portdistribute/api/order/refundTicket", AuthRole: "user"},
	"refund-product": {Path: "/portdistribute/api/order/refundProduct", AuthRole: "user"},
	"balance":        {Path: "/portdistribute/api/preDeposit/searchBalance", AuthRole: "admin", AllowEmpty: true},
	"change-fee":     {Path: "/portdistribute/api/order/searchChangeFee", AuthRole: "user"},
	"change-voyages": {Path: "/portdistribute/api/voyage/searchVoyageForChange", AuthRole: "user"},
	"change-lock":    {Path: "/portdistribute/api/order/changeOrderLock", AuthRole: "user"},
	"change-unlock":  {Path: "/portdistribute/api/order/changeOrderUnLock", AuthRole: "user"},
}

func CruiseProxy(action string) gin.HandlerFunc {
	endpoint, ok := cruiseEndpoints[action]
	return func(c *gin.Context) {
		if !ok {
			c.JSON(http.StatusNotFound, dto.Fail(404, "游船接口不存在"))
			return
		}
		if endpoint.AuthRole != "" {
			role, _ := c.Get("actor_role")
			if role != endpoint.AuthRole {
				c.JSON(http.StatusForbidden, dto.Fail(403, "无权调用该游船接口"))
				return
			}
		}

		body := map[string]interface{}{}
		if c.Request.Body != nil && c.Request.ContentLength != 0 {
			if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
				c.JSON(http.StatusBadRequest, dto.Fail(400, "请求 JSON 格式错误"))
				return
			}
		}
		if !endpoint.AllowEmpty && len(body) == 0 {
			c.JSON(http.StatusBadRequest, dto.Fail(400, "请求参数不能为空"))
			return
		}

		client, err := jzg.NewClient(config.Load().JZGBaseURL, config.Load().JZGDistributorCode, config.Load().JZGAccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Fail(500, err.Error()))
			return
		}

		resp, err := client.Post(c.Request.Context(), endpoint.Path, body)
		if err != nil {
			c.JSON(http.StatusBadGateway, dto.Fail(502, err.Error()))
			return
		}
		if resp.Code != "200" || !resp.Success {
			c.JSON(http.StatusBadGateway, dto.Fail(502, "九洲港接口返回失败: "+resp.Msg))
			return
		}

		var data interface{}
		if len(resp.Data) > 0 && string(resp.Data) != "null" {
			if err := json.Unmarshal(resp.Data, &data); err != nil {
				data = json.RawMessage(resp.Data)
			}
		}
		c.JSON(http.StatusOK, dto.Success(data))
	}
}
