package domain

import "time"

type OutboxEvent struct {
	ID             OutboxID
	NotificationID NotificationID
	EventType      string // "notification.created"
	Payload        []byte // serialized event
	Published      bool
	PublishedAt    *time.Time
	CreatedAt      time.Time
}

const EventTypeNotificationCreated = "notification.created"
const EventTypePreferenceUpdated = "user.preference.updated"
