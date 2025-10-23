package pagination

type GetListRequest struct {
	Page    int
	Size    int `json:"size"`
	SortBy  string
	SortDir string
	Query   string
	Filter  []FilterItem
}

// Implement GetPage method
func (p *GetListRequest) GetPage() int {
	return p.Page
}

// Implement GetSize method
func (p *GetListRequest) GetSize() int {
	return p.Size
}

// Implement GetSortBy method
func (p *GetListRequest) GetSortBy() string {
	return p.SortBy
}

// Implement GetSortDir method
func (p *GetListRequest) GetSortDir() string {
	return p.SortDir
}

// Implement GetQuery method
func (p *GetListRequest) GetQuery() string {
	return p.Query
}

// Implement GetFilter method
func (p *GetListRequest) GetFilter() []FilterItem {
	return p.Filter
}

// Implement GetFilter method
func (p *GetListRequest) AddFilter(f FilterItem) {
	if p.Filter == nil {
		// p.Filter = new([]FilterItem)
		p.Filter = make([]FilterItem, 0)
	}
	p.Filter = append(p.Filter, f)
}

func (p *GetListRequest) GenerateFilter() {
	// do nothing
}
