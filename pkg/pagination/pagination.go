package pagination

// Cursor holds cursor entity which is used to create listed page.
type Cursor struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

// Page holds page entity which is used to show listed page.
type Page struct {
	Offset      int `json:"offset"`
	Limit       int `json:"limit"`
	CurrentPage int `json:"currentPage"`
	PageCount   int `json:"pageCount"`
	TotalCount  int `json:"totalCount"`
}
