package template

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, req CreateTemplateRequest) (*TemplateResponse, error) {
	variables := req.Variables
	if variables == nil {
		variables = extractVariables(req.Subject, req.Body)
	}

	t := Template{
		ID:        uuid.New().String(),
		Version:   1,
		Name:      req.Name,
		Channel:   req.Channel,
		Subject:   req.Subject,
		Body:      req.Body,
		Variables: variables,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.store.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create template: %w", err)
	}

	return toResponse(t), nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdateTemplateRequest) (*TemplateResponse, error) {
	current, err := s.store.GetActive(ctx, id)
	if err != nil {
		return nil, err
	}

	variables := req.Variables
	if variables == nil {
		variables = extractVariables(req.Subject, req.Body)
	}

	next := current
	next.Subject = req.Subject
	next.Body = req.Body
	next.Variables = variables
	next.Version = current.Version + 1
	next.IsActive = true
	next.CreatedAt = time.Now().UTC()

	if err := s.store.CreateNewVersion(ctx, next); err != nil {
		return nil, fmt.Errorf("update template: %w", err)
	}

	return toResponse(next), nil
}

func (s *Service) GetActive(ctx context.Context, id string) (*TemplateResponse, error) {
	t, err := s.store.GetActive(ctx, id)
	if err != nil {
		return nil, err
	}
	return toResponse(t), nil
}

func (s *Service) GetVersion(ctx context.Context, id string, version int) (*TemplateResponse, error) {
	t, err := s.store.GetVersion(ctx, id, version)
	if err != nil {
		return nil, err
	}
	return toResponse(t), nil
}

func (s *Service) List(ctx context.Context, limit int, cursor string) (*ListTemplatesResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	templates, next, err := s.store.List(ctx, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}

	resp := &ListTemplatesResponse{
		Templates: make([]TemplateResponse, len(templates)),
		NextCursor: next,
	}
	for i, t := range templates {
		resp.Templates[i] = *toResponse(t)
	}
	return resp, nil
}

func (s *Service) Deactivate(ctx context.Context, id string) error {
	return s.store.Deactivate(ctx, id)
}

func toResponse(t Template) *TemplateResponse {
	return &TemplateResponse{
		ID:        t.ID,
		Version:   t.Version,
		Name:      t.Name,
		Channel:   t.Channel,
		Subject:   t.Subject,
		Body:      t.Body,
		Variables: t.Variables,
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
	}
}

func extractVariables(subject, body string) []string {
	seen := map[string]bool{}
	var vars []string

	for _, text := range []string{subject, body} {
		for {
			start := strings.Index(text, "{{")
			if start == -1 {
				break
			}
			end := strings.Index(text[start:], "}}")
			if end == -1 {
				break
			}
			name := strings.TrimSpace(text[start+2 : start+end])
			if name != "" && !seen[name] {
				seen[name] = true
				vars = append(vars, name)
			}
			text = text[start+end+2:]
		}
	}
	return vars
}
