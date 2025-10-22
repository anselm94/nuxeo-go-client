package nuxeo

type EntityType string

type PaginationOptions struct {
	CurrentPageIndex string `json:"currentPageIndex"`
	PageSize         string `json:"pageSize"`
}

type SortedPaginationOptions struct {
	CurrentPageIndex string `json:"currentPageIndex"`
	PageSize         string `json:"pageSize"`
	MaxResults       string `json:"maxResults"`
	SortBy           string `json:"sortBy"`
	SortOrder        string `json:"sortOrder"`
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
