package airtable

import (
	"context"
	"fmt"
)

type Record struct {
	ID          string         `json:"id,omitempty"`
	Fields      map[string]any `json:"fields"`
	CreatedTime string         `json:"createdTime,omitempty"`
	Deleted     bool           `json:"deleted,omitempty"`

	c *Client
	b *Base
	t *Table
}

type ListRecordOptions struct {
}

func (t *Table) ListRecords(ctx context.Context, opts *ListRecordOptions) ([]*Record, error) {
	return listAll[Record](ctx, t.c, fmt.Sprintf("%s/%s", t.b.ID, t.ID), nil, "records", func(r *Record) error {
		r.c = t.c
		r.b = t.b
		r.t = t
		return nil
	})
}

func (r Record) String() string {
	return r.Fields[r.t.Fields[0].Name].(string)
}
