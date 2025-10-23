package pagination

type ResultPagination struct {
	Data        interface{} `json:"contents"`
	Count       int64       `json:"totalElements"`
	Page        int         `json:"pageNumber"`
	RowsPerPage int         `json:"size"`
	TotalPages  int         `json:"totalPages"`
}

type Pagination struct {
	Page        int   `json:"page"`
	Count       int64 `json:"totalElements"`
	RowsPerPage int   `json:"size"`
}

type PaginationRequestDto interface {
	GetPage() int
	GetSize() int
	GetSortBy() string
	GetSortDir() string
	GetQuery() string
	GetFilter() []FilterItem
	AddFilter(FilterItem)
	GenerateFilter()
}

type FilterItem struct {
	Field string
	Op    string
	Val   interface{}
}
