package models

type FilterResponse[T any] struct {
	Hits        T   `json:"hits"`
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	TotalHits   int `json:"totalHits"`
	HitsPerPage int `json:"hitsPerPage"`
}
