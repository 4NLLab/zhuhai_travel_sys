package models

import "time"

type JZGCallbackLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	EventType string    `gorm:"size:40;index;not null" json:"event_type"`
	OrderNo   *string   `gorm:"size:64;index" json:"order_no"`
	Verified  bool      `gorm:"not null;default:false" json:"verified"`
	RequestIP *string   `gorm:"size:64" json:"request_ip"`
	Payload   string    `gorm:"type:json;not null" json:"payload"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

func (JZGCallbackLog) TableName() string { return "jzg_callback_logs" }
