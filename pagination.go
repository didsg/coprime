package coprime

import (
	"net/url"
	"strconv"
)

type PaginationParams struct {
	Limit  int
	Before string
	After  string
	Extra  map[string]string
}

type PrimePaginationParams struct {
	NextCursor    string `json:"next_cursor"`
	SortDirection string `json:"sort_direction"`
	HasNext       bool   `json:"has_next"`
}

func (p *PaginationParams) Encode(direction string) string {
	values := url.Values{}

	if p.Limit > 0 {
		values.Add("limit", strconv.Itoa(p.Limit))
	}
	if p.Before != "" && direction == "prev" {
		values.Add("before", p.Before)
	}
	if p.After != "" && direction == "next" {
		values.Add("after", p.After)
	}

	for k, v := range p.Extra {
		values.Add(k, v)
	}

	return values.Encode()
}

func (p *PaginationParams) Done(direction string) bool {
	if p.Before == "" && direction == "prev" {
		return true
	}

	if p.After == "" && direction == "next" {
		return true
	}

	return false
}

func (p *PaginationParams) AddExtraParam(key, value string) {
	if p.Extra == nil {
		p.Extra = make(map[string]string)
	}
	p.Extra[key] = value
}

func Encode(p PrimePaginationParams) string {
	values := url.Values{}

	if p.HasNext {
		values.Add("cursor", p.NextCursor)
	}

	return values.Encode()
}
