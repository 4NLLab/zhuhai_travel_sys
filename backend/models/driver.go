package models

import "time"

// -------------------- 司机模块 --------------------

type Driver struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	DriverNo       string    `gorm:"uniqueIndex;size:40;not null" json:"driver_no"`
	Name           string    `gorm:"size:80;not null" json:"name"`
	Phone          string    `gorm:"size:32;index;not null" json:"phone"`
	IdCardNo       *string   `gorm:"size:255" json:"id_card_no"`
	Status         string    `gorm:"size:24;default:active;index" json:"status"`
	CommissionRate float64   `gorm:"type:decimal(6,4);default:0.0800" json:"commission_rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Driver) TableName() string { return "drivers" }

type Vehicle struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	DriverID  uint64    `gorm:"index;not null" json:"driver_id"`
	PlateNo   string    `gorm:"uniqueIndex;size:32;not null" json:"plate_no"`
	Model     *string   `gorm:"size:80" json:"model"`
	Seats     *int      `json:"seats"`
	Status    string    `gorm:"size:24;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Vehicle) TableName() string { return "vehicles" }

type DriverQRCode struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	DriverID  uint64    `gorm:"index;not null" json:"driver_id"`
	VehicleID *uint64   `gorm:"index" json:"vehicle_id"`
	Code      string    `gorm:"uniqueIndex;size:80;not null" json:"code"`
	Scene     string    `gorm:"size:80;default:seat" json:"scene"`
	Status    string    `gorm:"size:24;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DriverQRCode) TableName() string { return "driver_qr_codes" }

type DriverCommission struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	DriverID        uint64     `gorm:"index;not null" json:"driver_id"`
	OrderID         uint64     `gorm:"index;not null" json:"order_id"`
	OrderItemID     uint64     `gorm:"not null" json:"order_item_id"`
	CommissionNo    string     `gorm:"uniqueIndex;size:64;not null" json:"commission_no"`
	BaseAmount      float64    `gorm:"type:decimal(10,2);not null" json:"base_amount"`
	Rate            float64    `gorm:"type:decimal(6,4);not null" json:"rate"`
	CommissionAmount float64   `gorm:"type:decimal(10,2);not null" json:"commission_amount"`
	Status          string     `gorm:"size:32;default:pending" json:"status"`
	SettledAt       *time.Time `json:"settled_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (DriverCommission) TableName() string { return "driver_commissions" }

type DriverWithdrawal struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	DriverID     uint64     `gorm:"index;not null" json:"driver_id"`
	WithdrawalNo string     `gorm:"uniqueIndex;size:64;not null" json:"withdrawal_no"`
	Amount       float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Channel      string     `gorm:"size:32;default:alipay" json:"channel"`
	Account      *string    `gorm:"size:80" json:"account"`
	RealName     *string    `gorm:"size:80" json:"real_name"`
	Status       string     `gorm:"size:32;default:pending" json:"status"`
	ProcessedBy  *uint64    `json:"processed_by"`
	ProcessedAt  *time.Time `json:"processed_at"`
	Remark       *string    `gorm:"size:255" json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (DriverWithdrawal) TableName() string { return "driver_withdrawals" }

// -------------------- 其他模块 --------------------

type InvoiceTitle struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	TitleType string    `gorm:"size:24;default:company" json:"title_type"`
	TitleName string    `gorm:"size:160;not null" json:"title_name"`
	TaxNo     *string   `gorm:"size:80" json:"tax_no"`
	Email     *string   `gorm:"size:120" json:"email"`
	IsDefault int8      `gorm:"default:0" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (InvoiceTitle) TableName() string { return "invoice_titles" }

type Invoice struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	UserID         uint64     `gorm:"index;not null" json:"user_id"`
	OrderID        uint64     `gorm:"index;not null" json:"order_id"`
	InvoiceTitleID *uint64    `json:"invoice_title_id"`
	InvoiceNo      *string    `gorm:"size:80" json:"invoice_no"`
	Amount         float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status         string     `gorm:"size:32;default:requested" json:"status"`
	IssuedAt       *time.Time `json:"issued_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (Invoice) TableName() string { return "invoices" }

type UserFavorite struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:uk_favorites_user_product;not null" json:"user_id"`
	ProductID uint64    `gorm:"uniqueIndex:uk_favorites_user_product;index;not null" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserFavorite) TableName() string { return "user_favorites" }

type SupportTicket struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	OrderID   *uint64   `gorm:"index" json:"order_id"`
	Type      string    `gorm:"size:32;not null" json:"type"`
	Title     string    `gorm:"size:160;not null" json:"title"`
	Content   *string   `gorm:"type:text" json:"content"`
	Status    string    `gorm:"size:32;default:open;index" json:"status"`
	HandledBy *uint64   `json:"handled_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SupportTicket) TableName() string { return "support_tickets" }

type AuditLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ActorType string    `gorm:"size:24;index:idx_audit_actor;not null" json:"actor_type"`
	ActorID   *uint64   `gorm:"index:idx_audit_actor" json:"actor_id"`
	Action    string    `gorm:"size:80;not null" json:"action"`
	TargetType *string  `gorm:"size:80;index:idx_audit_target" json:"target_type"`
	TargetID  *uint64   `gorm:"index:idx_audit_target" json:"target_id"`
	IP        *string   `gorm:"size:64" json:"ip"`
	UserAgent *string   `gorm:"size:500" json:"user_agent"`
	Payload   *string   `gorm:"type:json" json:"payload"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
