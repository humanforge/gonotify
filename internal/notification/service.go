package notification

import (
	"context"
	"fmt"
)

type notificationStore interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string) ([]Notification, error)
	MarkRead(ctx context.Context, id string) error
}

type Service struct {
	store notificationStore
}

func NewService(store notificationStore) *Service {
	return &Service{store: store}
}

func (s *Service) Send(ctx context.Context, userID, title, body string) error {
	n := &Notification{
		UserID: userID,
		Title:  title,
		Body:   body,
	}
	if err := s.store.Create(ctx, n); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*Notification, error) {
	n, err := s.store.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return n, nil
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]Notification, error) {
	notifications, err := s.store.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	return notifications, nil
}

func (s *Service) MarkRead(ctx context.Context, id string) error {
	if err := s.store.MarkRead(ctx, id); err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	return nil
}
