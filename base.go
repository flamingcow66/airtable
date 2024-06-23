package airtable

import (
	"context"
	"fmt"
)

type Base struct {
	ID              *string `json:"id,omitempty"`
	Name            *string `json:"name,omitempty"`
	PermissionLevel *string `json:"permissionLevel,omitempty"`

	c *Client
}

func (c *Client) ListBases(ctx context.Context) ([]*Base, error) {
	return listAll[Base](ctx, c, "meta/bases", nil, "bases", func(b *Base) error {
		b.c = c
		return nil
	})
}

func (c *Client) GetBaseByName(ctx context.Context, name string) (*Base, error) {
	bases, err := c.ListBases(ctx)
	if err != nil {
		return nil, err
	}

	for _, base := range bases {
		if *base.Name == name {
			return base, nil
		}
	}

	return nil, fmt.Errorf("base '%s' not found", name)
}

func (b Base) String() string {
	return *b.Name
}
