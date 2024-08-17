package utils

type GenericResponse struct {
	Msg    string                 `json:"msg,omitempty"`
	Status int                    `json:"status,omitempty"`
	Data   interface{}            `json:"data,omitempty"`
	Errors map[string]interface{} `json:"error,omitempty"`

	*PaginationResponse
}

type PaginationResponse struct {
	Count   int         `json:"count"`
	Results interface{} `json:"results"`
	Offset  int64       `json:"offset"`
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}
