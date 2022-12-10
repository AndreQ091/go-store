package frontend

import (
	"fmt"
	"go-store/core"
)

type Frontend interface {
	Start(kv *core.KeyValueStore) error
}

func NewFrontend(frontend string) (Frontend, error) {
	switch frontend {
	case "rest":
		return &restFrontend{}, nil
	default:
		return nil, fmt.Errorf("no such frontend %s", frontend)
	}
}
