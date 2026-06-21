package handlers

import (
	"net/http"

	"zhuhai_travel_backend/database"
	"zhuhai_travel_backend/dto"
	"zhuhai_travel_backend/models"

	"github.com/gin-gonic/gin"
)

// ProductList 产品列表（支持分类筛选、搜索、分页）
func ProductList(c *gin.Context) {
	categoryID := c.Query("category_id")
	keyword := c.Query("keyword")
	productType := c.Query("type")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	q := database.DB.Model(&models.Product{}).Where("status = ?", "active")
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q = q.Where("title LIKE ? OR subtitle LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if productType != "" {
		q = q.Where("product_type = ?", productType)
	}

	var total int64
	q.Count(&total)

	var products []models.Product
	q.Preload("SKUs").Preload("Images").
		Order("sort_order, id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&products)

	c.JSON(http.StatusOK, dto.Page(products, total, page, size))
}

// ProductDetail 产品详情
func ProductDetail(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.Preload("SKUs").Preload("Images").Preload("SKUs.Schedules").
		First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Fail(404, "产品不存在"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(product))
}

// ProductSchedulesQuery 查询某产品/具体日期的排期
func ProductSchedulesQuery(c *gin.Context) {
	productID := c.Query("product_id")
	date := c.Query("date") // 格式: 2026-07-01

	if productID == "" {
		c.JSON(http.StatusBadRequest, dto.Fail(400, "缺少 product_id"))
		return
	}

	q := database.DB.Model(&models.ProductSchedule{}).
		Where("product_id = ? AND status = ?", productID, "active")

	if date != "" {
		q = q.Where("travel_date = ?", date)
	}

	var schedules []models.ProductSchedule
	q.Order("travel_date, start_time").Find(&schedules)

	c.JSON(http.StatusOK, dto.Success(schedules))
}

// BannerList 轮播图列表
func BannerList(c *gin.Context) {
	var banners []models.Banner
	database.DB.Where("status = ?", "active").Order("sort_order").Find(&banners)
	c.JSON(http.StatusOK, dto.Success(banners))
}

// CategoryList 分类列表
func CategoryList(c *gin.Context) {
	var categories []models.ProductCategory
	database.DB.Where("parent_id IS NULL AND status = ?", "active").
		Preload("Children", "status = ?", "active").
		Order("sort_order").Find(&categories)
	c.JSON(http.StatusOK, dto.Success(categories))
}
