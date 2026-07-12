package notification

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, n *Notification) error {
	if err := s.db.WithContext(ctx).Create(n).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (s *Store) FindByID(ctx context.Context, id string) (*Notification, error) {
	var n Notification
	if err := s.db.WithContext(ctx).First(&n, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find notification by id: %w", err)
	}
	return &n, nil
}

func (s *Store) ListByUser(ctx context.Context, userID string) ([]Notification, error) {
	var notifications []Notification
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("list notifications for user: %w", err)
	}
	return notifications, nil
}

func (s *Store) MarkRead(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Model(&Notification{}).Where("id = ?", id).Update("read", true).Error; err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}
	return nil
}
