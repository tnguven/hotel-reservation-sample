package types

type (
	ResGeneric struct {
		Msg    string                 `json:"msg,omitempty"`
		Status int                    `json:"status,omitempty"`
		Data   interface{}            `json:"data,omitempty"`
		Errors map[string]interface{} `json:"error,omitempty"`
	}

	ResWithPaginate[T any] struct {
		ResGeneric
		Pagination *T
	}

	ResNumericPaginate struct {
		Count int64 `json:"count,omitempty"`
		Page  int   `json:"page,omitempty"`
		Limit int   `json:"limit,omitempty"`
	}

	ResCursorPaginate struct {
		LastID string `json:"lastID,omitempty"`
		Limit  int    `json:"limit"`
		Count  int64  `json:"count"`
	}
)
