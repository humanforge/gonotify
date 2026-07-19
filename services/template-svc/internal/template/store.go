package template

import (
	"context"
)

type Store interface {
	Create(ctx context.Context, t Template) error
	GetActive(ctx context.Context, id string) (Template, error)
	GetVersion(ctx context.Context, id string, version int) (Template, error)
	CreateNewVersion(ctx context.Context, next Template) error
	Deactivate(ctx context.Context, id string) error
	List(ctx context.Context, limit int, cursor string) ([]Template, string, error)
}
