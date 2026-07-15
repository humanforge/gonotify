package domain

import "time"

type DeliveryRecord struct {
	ID                DeliveryID
	NotificationID    NotificationID
	Channel           Channel
	Provider          Provider
	ProviderMessageID string
	ProviderResponse  map[string]any
	Status            DeliveryState
	SentAt            *time.Time
	DeliveredAt       *time.Time
	ErrorMessage      string
	RetryCount        int
}

func (d DeliveryRecord) IsTerminal() bool {
	return d.Status == DeliveryDelivered || d.Status == DeliveryFailed
}
