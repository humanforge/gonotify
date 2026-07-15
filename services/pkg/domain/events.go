package domain

import "time"

// NotificationCreatedEvent is what gets serialized into the outbox payload
// and published to Kafka. Kept separate from Notification because the
// wire contract should evolve independently of the DB schema.
type NotificationCreatedEvent struct {
	NotificationID NotificationID      `json:"notification_id"`
	Channel        Channel             `json:"channel"`
	Priority       Priority            `json:"priority"`
	Payload        NotificationPayload `json:"payload"`
	RetryCount     int                 `json:"retry_count"`
}

type PreferenceUpdatedEvent struct {
	ClientID       ClientID `json:"client_id"`
	ExternalUserID string   `json:"external_user_id"`
}

type WebhookDeliveryEvent struct {
	Provider          Provider      `json:"provider"`
	ProviderMessageID string        `json:"provider_message_id"`
	Status            DeliveryState `json:"status"`
	RawPayload        []byte        `json:"raw_payload"`
	ReceivedAt        time.Time     `json:"received_at"`
}
