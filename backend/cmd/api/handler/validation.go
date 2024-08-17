package handler

type PaginationFilter struct {
	Limit  int64 `validate:"numeric,max=60,omitempty"`
	Offset int64 `validate:"numeric,omitempty"`
}
