package models

import "time"

// -------------------- 用户模块 --------------------

type User struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	Openid         *string    `gorm:"uniqueIndex;size:80" json:"openid"`
	Unionid        *string    `gorm:"uniqueIndex;size:80" json:"unionid"`
	Phone          *string    `gorm:"size:32;index" json:"phone"`
	Nickname       *string    `gorm:"size:80" json:"nickname"`
	AvatarURL      *string    `gorm:"size:500" json:"avatar_url"`
	RealName       *string    `gorm:"size:80" json:"real_name"`
	IdCardNo       *string    `gorm:"size:255" json:"id_card_no"` // VARBINARY
	RealnameStatus string     `gorm:"size:24;default:unverified" json:"realname_status"`
	Status         string     `gorm:"size:24;default:active;index" json:"status"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (User) TableName() string { return "users" }

type Traveler struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"size:80;not null" json:"name"`
	Phone     *string   `gorm:"size:32" json:"phone"`
	IdType    string    `gorm:"size:24;default:id_card" json:"id_type"`
	IdNo      string    `gorm:"size:255;not null" json:"id_no"`
	IsDefault int8      `gorm:"default:0" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Traveler) TableName() string { return "travelers" }

type AdminUser struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:80;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	DisplayName  string     `gorm:"size:80;not null" json:"display_name"`
	Role         string     `gorm:"size:32;default:operator" json:"role"`
	Status       string     `gorm:"size:24;default:active" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (AdminUser) TableName() string { return "admin_users" }

// -------------------- 产品模块 --------------------

type ProductCategory struct {
	ID        uint64            `gorm:"primaryKey" json:"id"`
	ParentID  *uint64           `gorm:"index" json:"parent_id"`
	Name      string            `gorm:"size:80;not null" json:"name"`
	Slug      string            `gorm:"uniqueIndex;size:80;not null" json:"slug"`
	Icon      *string           `gorm:"size:80" json:"icon"`
	SortOrder int               `gorm:"default:0" json:"sort_order"`
	Status    string            `gorm:"size:24;default:active" json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Children  []ProductCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (ProductCategory) TableName() string { return "product_categories" }

type Product struct {
	ID           uint64          `gorm:"primaryKey" json:"id"`
	CategoryID   uint64          `gorm:"index;not null" json:"category_id"`
	Title        string          `gorm:"size:160;not null" json:"title"`
	Subtitle     *string         `gorm:"size:255" json:"subtitle"`
	ProductType  string          `gorm:"size:32;not null" json:"product_type"`
	CoverURL     *string         `gorm:"size:500" json:"cover_url"`
	LocationName *string         `gorm:"size:120" json:"location_name"`
	Address      *string         `gorm:"size:255" json:"address"`
	Notice       *string         `gorm:"type:text" json:"notice"`
	RefundPolicy *string         `gorm:"type:text" json:"refund_policy"`
	Status       string          `gorm:"size:24;default:draft" json:"status"`
	SortOrder    int             `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Category     ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SKUs         []ProductSKU    `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Images       []ProductImage  `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

func (Product) TableName() string { return "products" }

type ProductImage struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ProductID uint64    `gorm:"index;not null" json:"product_id"`
	ImageURL  string    `gorm:"size:500;not null" json:"image_url"`
	AltText   *string   `gorm:"size:160" json:"alt_text"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (ProductImage) TableName() string { return "product_images" }

type ProductSKU struct {
	ID              uint64            `gorm:"primaryKey" json:"id"`
	ProductID       uint64            `gorm:"index;not null" json:"product_id"`
	SkuName         string            `gorm:"size:120;not null" json:"sku_name"`
	MarketPrice     *float64          `gorm:"type:decimal(10,2)" json:"market_price"`
	SalePrice       float64           `gorm:"type:decimal(10,2);not null" json:"sale_price"`
	SettlementPrice *float64          `gorm:"type:decimal(10,2)" json:"settlement_price"`
	StockMode       string            `gorm:"size:24;default:schedule" json:"stock_mode"`
	Status          string            `gorm:"size:24;default:active;index" json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Schedules       []ProductSchedule `gorm:"foreignKey:SkuID" json:"schedules,omitempty"`
}

func (ProductSKU) TableName() string { return "product_skus" }

type ProductSchedule struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	ProductID   uint64    `gorm:"index;not null" json:"product_id"`
	SkuID       uint64    `gorm:"uniqueIndex:uk_schedule_sku_time;not null" json:"sku_id"`
	TravelDate  string    `gorm:"type:date;not null;uniqueIndex:uk_schedule_sku_time" json:"travel_date"`
	StartTime   *string   `gorm:"type:time;uniqueIndex:uk_schedule_sku_time" json:"start_time"`
	EndTime     *string   `gorm:"type:time" json:"end_time"`
	Venue       *string   `gorm:"size:120" json:"venue"`
	TotalStock  int       `gorm:"not null;default:0" json:"total_stock"`
	LockedStock int       `gorm:"not null;default:0" json:"locked_stock"`
	SoldStock   int       `gorm:"not null;default:0" json:"sold_stock"`
	Status      string    `gorm:"size:24;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ProductSchedule) TableName() string { return "product_schedules" }

type Banner struct {
	ID         uint64     `gorm:"primaryKey" json:"id"`
	Title      string     `gorm:"size:120;not null" json:"title"`
	Subtitle   *string    `gorm:"size:255" json:"subtitle"`
	ImageURL   string     `gorm:"size:500;not null" json:"image_url"`
	LinkType   string     `gorm:"size:32;default:product" json:"link_type"`
	LinkTarget *string    `gorm:"size:160" json:"link_target"`
	SortOrder  int        `gorm:"default:0" json:"sort_order"`
	Status     string     `gorm:"size:24;default:active" json:"status"`
	StartsAt   *time.Time `json:"starts_at"`
	EndsAt     *time.Time `json:"ends_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Banner) TableName() string { return "banners" }
