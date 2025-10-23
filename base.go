package nuxeo

import (
	"fmt"
	"net/url"
)

type EntityType string

type PaginationOptions struct {
	CurrentPageIndex int `json:"currentPageIndex"`
	PageSize         int `json:"pageSize"`
}

func (p *PaginationOptions) QueryParams() string {
	if p == nil {
		return ""
	}
	queryParams := url.Values{}
	if p.CurrentPageIndex > -1 {
		queryParams.Add("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Add("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	return queryParams.Encode()
}

type SortedPaginationOptions struct {
	CurrentPageIndex int    `json:"currentPageIndex"`
	PageSize         int    `json:"pageSize"`
	MaxResults       int    `json:"maxResults"`
	SortBy           string `json:"sortBy"`
	SortOrder        string `json:"sortOrder"`
}

func (p *SortedPaginationOptions) QueryParams() string {
	if p == nil {
		return ""
	}
	queryParams := url.Values{}
	if p.CurrentPageIndex > -1 {
		queryParams.Add("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Add("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	if p.MaxResults != 0 {
		queryParams.Add("maxResults", fmt.Sprintf("%d", p.MaxResults))
	}
	if p.SortBy != "" {
		queryParams.Add("sortBy", p.SortBy)
	}
	if p.SortOrder != "" {
		queryParams.Add("sortOrder", p.SortOrder)
	}
	return queryParams.Encode()
}

type entity struct {
	EntityType EntityType `json:"entity-type"`
}

type entities[T any] struct {
	entity
	Entries           []T            `json:"entries"`
	ContextParameters map[string]any `json:"contextParameters"`
}

type paginableEntities[T any] struct {
	entity
	IsPaginable             bool `json:"isPaginable"`
	ResultsCount            int  `json:"resultsCount"`
	PageSize                int  `json:"pageSize"`
	MaxPageSize             int  `json:"maxPageSize"`
	ResultsCountLimit       int  `json:"resultsCountLimit"`
	CurrentPageSize         int  `json:"currentPageSize"`
	CurrentPageIndex        int  `json:"currentPageIndex"`
	CurrentPageOffset       int  `json:"currentPageOffset"`
	NumberOfPages           int  `json:"numberOfPages"`
	IsPreviousPageAvailable bool `json:"isPreviousPageAvailable"`
	IsNextPageAvailable     bool `json:"isNextPageAvailable"`
	IsLastPageAvailable     bool `json:"isLastPageAvailable"`
	IsSortable              bool `json:"isSortable"`
	HasError                bool `json:"hasError"`
	ErrorMessage            any  `json:"errorMessage"`
	PageIndex               int  `json:"pageIndex"`
	PageCount               int  `json:"pageCount"`
	Entries                 []T  `json:"entries"`
}
