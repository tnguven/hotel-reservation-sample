package utils

type GenericResponse struct {
	Msg    string                 `json:"msg,omitempty"`
	Status int                    `json:"status,omitempty"`
	Data   interface{}            `json:"data,omitempty"`
	Errors map[string]interface{} `json:"error,omitempty"`
}
