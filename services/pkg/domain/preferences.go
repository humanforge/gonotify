// Package domain contains the core business entities, value objects,
// and business rules for the notification service.
package domain

import "time"

type Preferences struct {
	ClientID       ClientID
	ExternalUserID ExternalUserID
	Channels       map[Channel]bool
	Types          map[NotificationType]bool
	DoNotDisturb   bool
	UpdatedAt      time.Time
}

// Allows reports whether the notification is permitted for the given
// channel and notification type based on the user's preferences.
// It returns false if Do Not Disturb is enabled, if the channel is
// disabled, or if the notification type is disabled.
func (p Preferences) Allows(ch Channel, t NotificationType) bool {
	if p.DoNotDisturb {
		return false
	}
	if enabled, ok := p.Channels[ch]; ok && !enabled {
		return false
	}
	if enabled, ok := p.Types[t]; ok && !enabled {
		return false
	}
	return true
}

// FilterChannels returns only the channels the user has opted into,
// given the requested channels and notification type.
func (p Preferences) FilterChannels(requested []Channel, t NotificationType) []Channel {
	var allowed []Channel
	for _, ch := range requested {
		if p.Allows(ch, t) {
			allowed = append(allowed, ch)
		}
	}
	return allowed
}
