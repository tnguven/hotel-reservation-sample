package handler

type PaginationFilter struct {
	Limit int64 `validate:"required,numeric,max=60"`
	Page  int64 `validate:"required,numeric"`
}
