package template

import "errors"

var (
	ErrTemplateNotFound = errors.New("template not found")
	ErrTemplateInactive = errors.New("template is not active")
	ErrConcurrentUpdate = errors.New("concurrent update detected, retry")
)
