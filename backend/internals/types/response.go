package types

type GenericResponse struct {
	Msg    string                 `json:"msg,omitempty"`
	Status int                    `json:"status,omitempty"`
	Data   interface{}            `json:"data,omitempty"`
	Errors map[string]interface{} `json:"error,omitempty"`

	*PaginationResponse
}

type PaginationResponse struct {
	Count int64 `json:"count"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
