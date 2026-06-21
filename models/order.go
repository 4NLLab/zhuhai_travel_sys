package models

import (
	"time"
)

// -------------------- 订单模块 --------------------

type Order struct {
	ID             uint64      `gorm:"primaryKey" json:"id"`
	OrderNo        string      `gorm:"uniqueIndex;size:40;not null" json:"order_no"`
	UserID         uint64      `gorm:"index;not null" json:"user_id"`
	Source         string      `gorm:"size:32;default:miniapp" json:"source"`
	DriverID       *uint64     `gorm:"index" json:"driver_id"`
	DriverQRCodeID *uint64     `json:"driver_qr_code_id"`
	Status         string      `gorm:"size:32;default:pending_payment;index" json:"status"`
	TotalAmount    float64     `gorm:"type:decimal(10,2);default:0" json:"total_amount"`
	DiscountAmount float64     `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	PayableAmount  float64     `gorm:"type:decimal(10,2);default:0" json:"payable_amount"`
	PaidAmount     float64     `gorm:"type:decimal(10,2);default:0" json:"paid_amount"`
	ContactName    *string     `gorm:"size:80" json:"contact_name"`
	ContactPhone   *string     `gorm:"size:32" json:"contact_phone"`
	Remark         *string     `gorm:"size:255" json:"remark"`
	PaidAt         *time.Time  `json:"paid_at"`
	CancelledAt    *time.Time  `json:"cancelled_at"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payment        *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	ID           uint64          `gorm:"primaryKey" json:"id"`
	OrderID      uint64          `gorm:"index;not null" json:"order_id"`
	ProductID    uint64          `gorm:"index;not null" json:"product_id"`
	SkuID        uint64          `gorm:"not null" json:"sku_id"`
	ScheduleID   *uint64         `gorm:"index" json:"schedule_id"`
	ProductTitle string          `gorm:"size:160;not null" json:"product_title"`
	SkuName      string          `gorm:"size:120;not null" json:"sku_name"`
	TravelDate   *string         `gorm:"type:date" json:"travel_date"`
	StartTime    *string         `gorm:"type:time" json:"start_time"`
	Quantity     uint            `gorm:"not null" json:"quantity"`
	UnitPrice    float64         `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice   float64         `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status       string          `gorm:"size:32;default:active" json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	Tickets      []Ticket        `gorm:"foreignKey:OrderItemID" json:"tickets,omitempty"`
	Travelers    []OrderTraveler `gorm:"foreignKey:OrderItemID" json:"travelers,omitempty"`
}

func (OrderItem) TableName() string { return "order_items" }

type OrderTraveler struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	OrderItemID uint64    `gorm:"index;not null" json:"order_item_id"`
	TravelerID  *uint64   `gorm:"index" json:"traveler_id"`
	Name        string    `gorm:"size:80;not null" json:"name"`
	Phone       *string   `gorm:"size:32" json:"phone"`
	IdType      string    `gorm:"size:24;default:id_card" json:"id_type"`
	IdNo        string    `gorm:"size:255;not null" json:"id_no"`
	CreatedAt   time.Time `json:"created_at"`
}

func (OrderTraveler) TableName() string { return "order_travelers" }

type Payment struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	OrderID       uint64     `gorm:"uniqueIndex;index;not null" json:"order_id"`
	PaymentNo     string     `gorm:"uniqueIndex;size:64;not null" json:"payment_no"`
	Channel       string     `gorm:"size:32;default:wechat" json:"channel"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string     `gorm:"size:32;default:pending;index" json:"status"`
	TransactionID *string    `gorm:"size:128" json:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at"`
	RawPayload    *string    `gorm:"type:json" json:"raw_payload"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Payment) TableName() string { return "payments" }

// -------------------- 票务模块 --------------------

type Ticket struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	OrderItemID uint64     `gorm:"index;not null" json:"order_item_id"`
	TicketNo    string     `gorm:"uniqueIndex;size:64;not null" json:"ticket_no"`
	QRTokenHash string     `gorm:"uniqueIndex;size:64;not null" json:"qr_token_hash"`
	Status      string     `gorm:"size:32;default:valid;index" json:"status"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidTo     *time.Time `json:"valid_to"`
	UsedAt      *time.Time `json:"used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	OrderItem   OrderItem  `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
}

func (Ticket) TableName() string { return "tickets" }

type TicketVerification struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	TicketID        uint64    `gorm:"index;not null" json:"ticket_id"`
	VerifierAdminID *uint64   `json:"verifier_admin_id"`
	VerifyLocation  *string   `gorm:"size:120" json:"verify_location"`
	VerifyResult    string    `gorm:"size:32;not null" json:"verify_result"`
	VerifyNote      *string   `gorm:"size:255" json:"verify_note"`
	CreatedAt       time.Time `gorm:"index" json:"created_at"`
}

func (TicketVerification) TableName() string { return "ticket_verifications" }

// -------------------- 退款模块 --------------------

type Refund struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	OrderID     uint64     `gorm:"index;not null" json:"order_id"`
	OrderItemID *uint64    `json:"order_item_id"`
	RefundNo    string     `gorm:"uniqueIndex;size:64;not null" json:"refund_no"`
	Amount      float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Reason      *string    `gorm:"size:255" json:"reason"`
	Status      string     `gorm:"size:32;default:requested;index" json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	RawPayload  *string    `gorm:"type:json" json:"raw_payload"`
}

func (Refund) TableName() string { return "refunds" }
