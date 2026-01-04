package entities

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	OrderId   uint      `gorm:"index;not null"`
	Total     float32   `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
}

func (Payment) TableName() string {
	return "payment"
}
