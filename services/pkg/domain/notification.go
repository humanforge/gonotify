package domain

import "time"

type Notification struct {
	ID             NotificationID
	ClientID       ClientID
	ExternalUserID ExternalUserID
	TemplateID     TemplateID
	Channel        Channel
	Payload        any
	Status         NotificationStatus
	Priority       Priority
	Type           NotificationType
	ScheduledAt    *time.Time
	CreatedAt      time.Time
	LastUpdatedAt  time.Time
	Metadata       map[string]string
}

type NotificationPayload struct {
	Subject string
	Body    string
	To      string
	Extra   map[string]string
}

func (n Notification) IsScheduled() bool {
	return n.ScheduledAt != nil && n.ScheduledAt.After(time.Now())
}

func (n Notification) IsTerminal() bool {
	switch n.Status {
	case StatusDelivered, StatusFailed, StatusCancelled:
		return true
	}
	return false
}
