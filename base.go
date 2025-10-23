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

func (p *PaginationOptions) QueryParams() url.Values {
	if p == nil {
		return nil
	}
	queryParams := make(url.Values)
	if p.CurrentPageIndex > -1 {
		queryParams.Set("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Set("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	return queryParams
}

type SortedPaginationOptions struct {
	CurrentPageIndex int    `json:"currentPageIndex"`
	PageSize         int    `json:"pageSize"`
	MaxResults       int    `json:"maxResults"`
	SortBy           string `json:"sortBy"`
	SortOrder        string `json:"sortOrder"`
}

func (p *SortedPaginationOptions) QueryParams() url.Values {
	if p == nil {
		return nil
	}
	queryParams := make(url.Values)
	if p.CurrentPageIndex > -1 {
		queryParams.Set("currentPageIndex", fmt.Sprintf("%d", p.CurrentPageIndex))
	}
	if p.PageSize != 0 {
		queryParams.Set("pageSize", fmt.Sprintf("%d", p.PageSize))
	}
	if p.MaxResults != 0 {
		queryParams.Set("maxResults", fmt.Sprintf("%d", p.MaxResults))
	}
	if p.SortBy != "" {
		queryParams.Set("sortBy", p.SortBy)
	}
	if p.SortOrder != "" {
		queryParams.Set("sortOrder", p.SortOrder)
	}
	return queryParams
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
