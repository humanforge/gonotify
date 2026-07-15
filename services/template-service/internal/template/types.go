package template

import "time"

type Template struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Version   int       `gorm:"not null;default:1" json:"version"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Channel   string    `gorm:"type:varchar(20);not null;index" json:"channel"`
	Subject   string    `gorm:"type:varchar(255)" json:"subject"`
	Variables []string  `gorm:"serializer:json;type:jsonb" json:"variables"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}
