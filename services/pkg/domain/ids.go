package domain

import "github.com/google/uuid"

type NotificationID string
type TemplateID string
type OutboxID string
type DeliveryID string
type ExternalUserID string
type ClientID string

func NewNotificationID() NotificationID { return NotificationID(uuid.NewString()) }
func NewOutboxID() OutboxID             { return OutboxID(uuid.NewString()) }
func NewDeliveryID() DeliveryID         { return DeliveryID(uuid.NewString()) }
