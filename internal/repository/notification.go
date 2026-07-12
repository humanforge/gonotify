package repository

import (
	"context"
	"fmt"

	"notification-service/internal/domain"

	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	if err := r.db.WithContext(ctx).Create(n).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepo) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	var n domain.Notification
	if err := r.db.WithContext(ctx).First(&n, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("find notification by id: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepo) ListByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	var notifications []domain.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("list notifications for user: %w", err)
	}
	return notifications, nil
}

func (r *NotificationRepo) MarkRead(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ?", id).Update("read", true).Error; err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}
	return nil
}
