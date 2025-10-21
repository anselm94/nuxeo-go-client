package nuxeo

type PaginationOptions struct {
	CurrentPageIndex int `json:"currentPageIndex"`
	PageSize         int `json:"pageSize"`
}

type PaginatedEntities[T any] struct {
	EntityType              string `json:"entity-type"`
	IsPaginable             bool   `json:"isPaginable"`
	ResultsCount            int    `json:"resultsCount"`
	PageSize                int    `json:"pageSize"`
	MaxPageSize             int    `json:"maxPageSize"`
	ResultsCountLimit       int    `json:"resultsCountLimit"`
	CurrentPageSize         int    `json:"currentPageSize"`
	CurrentPageIndex        int    `json:"currentPageIndex"`
	CurrentPageOffset       int    `json:"currentPageOffset"`
	NumberOfPages           int    `json:"numberOfPages"`
	IsPreviousPageAvailable bool   `json:"isPreviousPageAvailable"`
	IsNextPageAvailable     bool   `json:"isNextPageAvailable"`
	IsLastPageAvailable     bool   `json:"isLastPageAvailable"`
	IsSortable              bool   `json:"isSortable"`
	HasError                bool   `json:"hasError"`
	ErrorMessage            any    `json:"errorMessage"`
	PageIndex               int    `json:"pageIndex"`
	PageCount               int    `json:"pageCount"`
	Entries                 []T    `json:"entries"`
}
