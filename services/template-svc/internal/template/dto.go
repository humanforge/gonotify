package template

import "time"

type CreateTemplateRequest struct {
	Name      string   `json:"name" validate:"required,min=3,max=255"`
	Channel   string   `json:"channel" validate:"required,oneof=EMAIL SMS PUSH INAPP"`
	Subject   string   `json:"subject,omitempty"`
	Body      string   `json:"body" validate:"required"`
	Variables []string `json:"variables,omitempty"`
}

type UpdateTemplateRequest struct {
	Subject   string   `json:"subject,omitempty"`
	Body      string   `json:"body" validate:"required"`
	Variables []string `json:"variables,omitempty"`
}

type TemplateResponse struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	Name      string    `json:"name"`
	Channel   string    `json:"channel"`
	Subject   string    `json:"subject,omitempty"`
	Body      string    `json:"body"`
	Variables []string  `json:"variables"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ListTemplatesRequest struct {
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Cursor string `json:"cursor,omitempty"`
}

type ListTemplatesResponse struct {
	Templates  []TemplateResponse `json:"templates"`
	NextCursor string             `json:"next_cursor,omitempty"`
}
