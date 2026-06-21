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

		// ==================== 用户端 ====================
		user := api.Group("/users")
		{
			user.GET("/profile", handlers.UserProfile)
			user.POST("/realname", handlers.UserRealnameSubmit)
		}

		// 出游人
		traveler := api.Group("/travelers")
		{
			traveler.GET("", handlers.TravelerList)
			traveler.POST("", handlers.TravelerCreate)
			traveler.PUT("/:id", handlers.TravelerUpdate)
			traveler.DELETE("/:id", handlers.TravelerDelete)
		}

		// 收藏
		fav := api.Group("/favorites")
		{
			fav.GET("", handlers.FavoriteList)
			fav.POST("/toggle", handlers.FavoriteToggle)
		}

		// 发票抬头
		invTitle := api.Group("/invoice-titles")
		{
			invTitle.GET("", handlers.InvoiceTitleList)
			invTitle.POST("", handlers.InvoiceTitleCreate)
			invTitle.DELETE("/:id", handlers.InvoiceTitleDelete)
		}

		// 发票
		api.POST("/invoices", handlers.InvoiceCreate)

		// ==================== 产品 ====================
		api.GET("/products", handlers.ProductList)
		api.GET("/products/schedules/query", handlers.ProductSchedulesQuery)
		api.GET("/products/:id", handlers.ProductDetail)
		api.GET("/categories", handlers.CategoryList)
		api.GET("/banners", handlers.BannerList)

		// ==================== 订单 ====================
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.OrderList)
		api.GET("/orders/:id", handlers.OrderDetail)
		api.POST("/payments/callback", handlers.PaymentCallback)

		// ==================== 票务 ====================
		ticket := api.Group("/tickets")
		{
			ticket.GET("/:id", handlers.TicketDetail)
			ticket.GET("/by-qr", handlers.TicketByQR)
			ticket.POST("/verify", handlers.TicketVerify)
			ticket.GET("/verifications", handlers.VerificationHistory)
		}

		// ==================== 司机端（小程序） ====================
		driverApp := api.Group("/driver")
		{
			driverApp.POST("/login", handlers.DriverLogin)
			driverApp.GET("/wallet", handlers.DriverWallet)
			driverApp.GET("/commissions", handlers.DriverViewCommissionList)
			driverApp.POST("/withdraw", handlers.DriverWithdraw)
			driverApp.GET("/withdrawals", handlers.DriverWithdrawalHistory)
		}

		// ==================== 管理端-司机 ====================
		driver := api.Group("/drivers")
		{
			driver.GET("", handlers.DriverList)
		}
		comm := api.Group("/commissions")
		{
			comm.GET("", handlers.DriverCommissionList)
			comm.GET("/summary", handlers.CommissionSummary)
		}

		// ==================== 管理后台 ====================
		admin := api.Group("/admin")
		{
			admin.POST("/login", handlers.AdminLogin)

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

			// 参数配置
			admin.GET("/params", handlers.AdminParams)
		}
	}

	return r
}
