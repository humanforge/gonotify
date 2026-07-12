package domain

import "context"

type Notification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string) ([]Notification, error)
	MarkRead(ctx context.Context, id string) error
}

type NotificationService interface {
	Send(ctx context.Context, userID, title, body string) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string) ([]Notification, error)
	MarkRead(ctx context.Context, id string) error
}
