package model

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Paging *PageMetaData `json:"paging,omitempty"`
	Errors string        `json:"errors,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data"`
	PageMetaData PageMetaData `json:"pageMetaData"`
}

type PageMetaData struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}
