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
		api.POST("/admin/login", handlers.AdminLogin)

		// ==================== 公开展示 ====================
		api.GET("/products", handlers.ProductList)
		api.GET("/products/schedules/query", handlers.ProductSchedulesQuery)
		api.GET("/products/:id", handlers.ProductDetail)
		api.GET("/categories", handlers.CategoryList)
		api.GET("/banners", handlers.BannerList)
		api.GET("/island-cruise/ports", handlers.IslandCruisePorts)
		api.GET("/island-cruise/voyages", handlers.IslandCruiseVoyages)
		api.GET("/island-cruise/price", handlers.IslandCruisePrice)

		// 支付回调由支付平台调用，使用独立签名校验。
		api.POST("/payments/callback", handlers.PaymentCallback)

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

		// ==================== 票务 ====================
		ticket := userArea.Group("/tickets")
		{
			ticket.GET("/by-qr", handlers.TicketByQR)
			ticket.GET("/verifications", handlers.VerificationHistory)
			ticket.GET("/:id", handlers.TicketDetail)
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
		{
			driver.GET("", handlers.DriverList)
		}
		comm := api.Group("/commissions")
		comm.Use(middleware.AuthRequired("admin"))
		{
			comm.GET("", handlers.DriverCommissionList)
			comm.GET("/summary", handlers.CommissionSummary)
		}

		// ==================== 管理后台 ====================
		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired("admin"))
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
			admin.POST("/drivers/review", handlers.AdminDriverReview)

			// 参数配置
			admin.GET("/params", handlers.AdminParams)
		}

		// 核销属于管理端能力，单独挂在 admin 权限下。
		api.POST("/tickets/verify", middleware.AuthRequired("admin"), handlers.TicketVerify)
	}

	return r
}
