package model

import "time"

type Session struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"index"`
	DeviceKey  string    `json:"device_key" gorm:"uniqueIndex:idx_user_device,size:64"`
	UserAgent  string    `json:"user_agent" gorm:"size:512"`
	IP         string    `json:"ip" gorm:"size:64"`
	LastActive time.Time `json:"last_active" gorm:"autoUpdateTime"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Status     int       `json:"status" gorm:"default:1"`
	Username   string    `json:"username" gorm:"-"`
}

const (
	SessionInactive = iota
	SessionActive
)

func (s *Session) IsActive() bool {
	return s.Status == SessionActive
}