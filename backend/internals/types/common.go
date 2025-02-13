package types

type PaginationQuery struct {
	Limit int64 `query:"limit" validate:"numeric,max=100,omitempty"`
	Page  int64 `query:"page" validate:"numeric,omitempty"`
}

func NewPaginationQuery(limit int, page int) *PaginationQuery {
	return &PaginationQuery{
		Limit: int64(limit),
		Page:  int64((page - 1) * limit),
	}
}
