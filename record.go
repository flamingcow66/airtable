package airtable

import (
	"context"
	"fmt"
)

type Record struct {
	ID          *string        `json:"id,omitempty"`
	Fields      map[string]any `json:"fields"`
	CreatedTime *string        `json:"createdTime,omitempty"`

	c *Client
	b *Base
	t *Table
}

type ListRecordOptions struct {
}

func (t *Table) ListRecords(ctx context.Context, opts *ListRecordOptions) ([]*Record, error) {
	return listAll[Record](ctx, t.c, fmt.Sprintf("%s/%s", *t.b.ID, *t.ID), nil, "records", func(r *Record) error {
		r.c = t.c
		r.b = t.b
		r.t = t
		return nil
	})
}

type updateRecordsRequest struct {
	PerformUpsert *performUpsertRequest `json:"performUpsert,omitempty"`
	Records       []*Record             `json:"records"`
}

type performUpsertRequest struct {
	FieldsToMergeOn []string `json:"fieldsToMergeOn"`
}

type updateRecordsResponse struct {
	Records []*Record `json:"records"`
}

func (t *Table) ReplaceRecords(ctx context.Context, records []*Record, matchFields []string) ([]*Record, error) {
	ret := []*Record{}

	req := &updateRecordsRequest{}

	if len(matchFields) > 0 {
		req.PerformUpsert = &performUpsertRequest{
			FieldsToMergeOn: matchFields,
		}
	}

	for _, chk := range chunk(records, 10) {
		req.Records = chk

		resp, err := put[updateRecordsResponse](ctx, t.c, fmt.Sprintf("%s/%s", *t.b.ID, *t.ID), req)
		if err != nil {
			return nil, err
		}

		for _, rec := range resp.Records {
			rec.c = t.c
			rec.b = t.b
			rec.t = t
			ret = append(ret, rec)
		}
	}

	return ret, nil
}

func (r Record) String() string {
	if r.ID == nil {
		return fmt.Sprintf("(%s [nil])", r.Fields[*r.t.Fields[0].Name].(string))
	} else {
		return fmt.Sprintf("(%s [%s])", r.Fields[*r.t.Fields[0].Name].(string), *r.ID)
	}
}
