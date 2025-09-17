package db

import (
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetSession(userID uint, deviceKey string) (*model.Session, error) {
	var s model.Session
	err := db.Where("user_id = ? AND device_key = ?", userID, deviceKey).First(&s).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find session")
	}
	return &s, nil
}

func CreateSession(s *model.Session) error {
	return errors.WithStack(db.Create(s).Error)
}

func UpsertSession(s *model.Session) error {
	return errors.WithStack(db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "device_key"}},
		UpdateAll: true,
	}).Create(s).Error)
}

func DeleteSession(userID uint, deviceKey string) error {
	return errors.WithStack(db.Where("user_id = ? AND device_key = ?", userID, deviceKey).Delete(&model.Session{}).Error)
}

func UpdateSessionLastActive(userID uint, deviceKey string) error {
	return errors.WithStack(db.Model(&model.Session{}).
		Where("user_id = ? AND device_key = ?", userID, deviceKey).
		Update("last_active", time.Now()).Error)
}

func DeactivateSession(deviceKey string) error {
	return errors.WithStack(db.Model(&model.Session{}).
		Where("device_key = ?", deviceKey).
		Update("status", model.SessionInactive).Error)
}

func ListUserSessions(userID uint, ipFilter string) ([]model.Session, error) {
	var sessions []model.Session
	query := db.Where("user_id = ? AND status = ?", userID, model.SessionActive)
	if ipFilter != "" {
		query = query.Where("ip LIKE ?", "%"+ipFilter+"%")
	}
	err := query.Order("last_active DESC").Find(&sessions).Error
	return sessions, errors.WithStack(err)
}

func ListSessions(usernameFilter, ipFilter string) ([]model.Session, error) {
	var sessions []model.Session
	query := db.Model(&model.Session{}).
		Select("sessions.*, users.username").
		Joins("LEFT JOIN users ON sessions.user_id = users.id").
		Where("sessions.status = ?", model.SessionActive)
	
	if usernameFilter != "" {
		query = query.Where("users.username LIKE ?", "%"+usernameFilter+"%")
	}
	if ipFilter != "" {
		query = query.Where("sessions.ip LIKE ?", "%"+ipFilter+"%")
	}
	
	err := query.Order("sessions.last_active DESC").Find(&sessions).Error
	return sessions, errors.WithStack(err)
}

func CleanExpiredSessions(expireBefore time.Time) error {
	return errors.WithStack(db.Where("last_active < ? AND status = ?", expireBefore, model.SessionActive).
		Update("status", model.SessionInactive).Error)
}