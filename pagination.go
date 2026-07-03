package moyasar

// PageParams contains common list pagination parameters.
//
// Moyasar list endpoints return 40 objects by default. If Page is omitted,
// Moyasar returns the first page containing the newest objects.
type PageParams struct {
	// Page is the requested page number.
	Page int
}

// PageMeta contains pagination metadata returned by list endpoints.
type PageMeta struct {
	// CurrentPage is the current page number.
	CurrentPage int `json:"current_page"`
	// NextPage is the next page number, or nil when there is no next page.
	NextPage *int `json:"next_page"`
	// PrevPage is the previous page number, or nil when there is no previous page.
	PrevPage *int `json:"prev_page"`
	// TotalPages is the total number of pages for the resource list.
	TotalPages int `json:"total_pages"`
	// TotalCount is the total number of objects when Moyasar includes it.
	TotalCount int `json:"total_count,omitempty"`
}
