package routes

import (
	"zhuhai_travel_backend/handlers"
	"zhuhai_travel_backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// ==================== 公开登录 ====================
		api.POST("/auth/phone-login", handlers.UserPhoneLogin)
		api.POST("/driver/register", handlers.DriverRegister)
		api.POST("/driver/login", handlers.DriverLogin)
		api.POST("/driver/upload-license", handlers.DriverUploadLicense)
		api.POST("/admin/login", handlers.AdminLogin)

		// ==================== 公开展示 ====================
		api.GET("/products", handlers.ProductList)
		api.GET("/products/schedules/query", handlers.ProductSchedulesQuery)
		api.GET("/products/:id", handlers.ProductDetail)
		api.GET("/categories", handlers.CategoryList)
		api.GET("/banners", handlers.BannerList)
		api.POST("/cruise/ports", handlers.CruiseProxy("ports"))
		api.POST("/cruise/cert-types", handlers.CruiseProxy("cert-types"))
		api.POST("/cruise/voyages", handlers.CruiseProxy("voyages"))
		api.POST("/cruise/prices", handlers.CruiseProxy("prices"))
		api.POST("/cruise/single-voyage", handlers.CruiseProxy("single-voyage"))

		// 支付回调由支付平台调用，使用独立签名校验。
		api.POST("/payments/callback", handlers.PaymentCallback)

		// 九洲港分销接口回调，由九洲港服务器主动通知。
		api.POST("/cruise/callback/refund-ticket", handlers.CruiseCallback("refund_ticket"))
		api.POST("/cruise/callback/refund-fee", handlers.CruiseCallback("refund_fee"))
		api.POST("/cruise/callback/check-ticket", handlers.CruiseCallback("check_ticket"))
		api.POST("/cruise/callback/refund-product", handlers.CruiseCallback("refund_product"))

		// ==================== 用户端 ====================
		userArea := api.Group("")
		userArea.Use(middleware.AuthRequired("user"))
		user := userArea.Group("/users")
		{
			user.GET("/profile", handlers.UserProfile)
			user.POST("/realname", handlers.UserRealnameSubmit)
		}

		// 出游人
		traveler := userArea.Group("/travelers")
		{
			traveler.GET("", handlers.TravelerList)
			traveler.POST("", handlers.TravelerCreate)
			traveler.PUT("/:id", handlers.TravelerUpdate)
			traveler.DELETE("/:id", handlers.TravelerDelete)
		}

		// 收藏
		fav := userArea.Group("/favorites")
		{
			fav.GET("", handlers.FavoriteList)
			fav.POST("/toggle", handlers.FavoriteToggle)
		}

		// 发票抬头
		invTitle := userArea.Group("/invoice-titles")
		{
			invTitle.GET("", handlers.InvoiceTitleList)
			invTitle.POST("", handlers.InvoiceTitleCreate)
			invTitle.DELETE("/:id", handlers.InvoiceTitleDelete)
		}

		// 发票
		userArea.POST("/invoices", handlers.InvoiceCreate)

		// ==================== 订单 ====================
		userArea.POST("/orders", handlers.CreateOrder)
		userArea.GET("/orders", handlers.OrderList)
		userArea.GET("/orders/:id", handlers.OrderDetail)
		userArea.POST("/cruise/lock-order", handlers.CruiseProxy("lock-order"))
		userArea.POST("/cruise/sale-order", handlers.CruiseProxy("sale-order"))
		userArea.POST("/cruise/unlock-order", handlers.CruiseProxy("unlock-order"))
		userArea.POST("/cruise/order", handlers.CruiseProxy("order"))
		userArea.POST("/cruise/refund-fee", handlers.CruiseProxy("refund-fee"))
		userArea.POST("/cruise/refund-ticket", handlers.CruiseProxy("refund-ticket"))
		userArea.POST("/cruise/refund-product", handlers.CruiseProxy("refund-product"))
		userArea.POST("/cruise/change-fee", handlers.CruiseProxy("change-fee"))
		userArea.POST("/cruise/change-voyages", handlers.CruiseProxy("change-voyages"))
		userArea.POST("/cruise/change-lock", handlers.CruiseProxy("change-lock"))
		userArea.POST("/cruise/change-unlock", handlers.CruiseProxy("change-unlock"))

		// ==================== 票务 ====================
		ticket := userArea.Group("/tickets")
		{
			ticket.GET("/:id", handlers.TicketDetail)
			ticket.GET("/by-qr", handlers.TicketByQR)
			ticket.GET("/verifications", handlers.VerificationHistory)
		}

		// ==================== 司机端（小程序） ====================
		driverApp := api.Group("/driver")
		driverApp.Use(middleware.AuthRequired("driver"))
		{
			driverApp.GET("/wallet", handlers.DriverWallet)
			driverApp.GET("/commissions", handlers.DriverViewCommissionList)
			driverApp.POST("/withdraw", handlers.DriverWithdraw)
			driverApp.GET("/withdrawals", handlers.DriverWithdrawalHistory)
		}

		// ==================== 管理端-司机 ====================
		driver := api.Group("/drivers")
		driver.Use(middleware.AuthRequired("admin"))
		driver.Use(middleware.AdminAuditLog())
		{
			driver.GET("", handlers.DriverList)
		}
		comm := api.Group("/commissions")
		comm.Use(middleware.AuthRequired("admin"))
		comm.Use(middleware.AdminAuditLog())
		{
			comm.GET("", handlers.DriverCommissionList)
			comm.GET("/summary", handlers.CommissionSummary)
		}

		// ==================== 管理后台 ====================
		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired("admin"))
		admin.Use(middleware.AdminAuditLog())
		{
			// 看板
			admin.GET("/dashboard", handlers.AdminDashboard)
			admin.GET("/trend", handlers.AdminTrend)

			// 订单管理
			admin.GET("/orders", handlers.AdminOrderList)
			admin.POST("/orders/refund", handlers.AdminOrderRefund)

			// 轮播图管理
			admin.GET("/banners", handlers.AdminBannerList)
			admin.POST("/banners", handlers.AdminBannerCreate)
			admin.PUT("/banners/:id", handlers.AdminBannerUpdate)
			admin.DELETE("/banners/:id", handlers.AdminBannerDelete)

			// 佣金结算
			admin.GET("/commission-batches", handlers.AdminCommissionBatches)
			admin.POST("/commissions/settle", handlers.AdminCommissionSettle)

			// 提现审核
			admin.GET("/withdrawals", handlers.AdminWithdrawalList)
			admin.POST("/withdrawals/process", handlers.AdminWithdrawalProcess)

			// 司机审核
			admin.POST("/drivers/review", handlers.AdminDriverReview)

			// 参数配置
			admin.GET("/params", handlers.AdminParams)

			// 操作日志
			admin.GET("/audit-logs", handlers.AdminAuditLogList)
			admin.GET("/audit-logs/:id", handlers.AdminAuditLogDetail)

			// 游船分销接口维护
			admin.POST("/cruise/balance", handlers.CruiseProxy("balance"))
		}

		api.POST("/tickets/verify", middleware.AuthRequired("admin"), middleware.AdminAuditLog(), handlers.TicketVerify)
	}

	return r
}
