package airtable

import (
	"context"
	"fmt"
)

type Table struct {
	ID             *string  `json:"id,omitempty"`
	PrimaryFieldID *string  `json:"primaryFieldId,omitempty"`
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Fields         []*Field `json:"fields"`
	Views          []*View  `json:"views"`

	c *Client
	b *Base
}

type Field struct {
	ID          *string        `json:"id"`
	Type        *string        `json:"type"`
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Options     map[string]any `json:"options"`
}

type View struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
	Name *string `json:"name"`
}

func (b *Base) ListTables(ctx context.Context) ([]*Table, error) {
	return listAll[Table](ctx, b.c, fmt.Sprintf("meta/bases/%s/tables", *b.ID), nil, "tables", func(t *Table) error {
		t.c = b.c
		t.b = b
		return nil
	})
}

func (b *Base) GetTableByName(ctx context.Context, name string) (*Table, error) {
	tables, err := b.ListTables(ctx)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		if *table.Name == name {
			return table, nil
		}
	}

	return nil, fmt.Errorf("table '%s' not found", name)
}

func (t Table) String() string {
	return fmt.Sprintf("%s.%s", t.b, t.Name)
}
