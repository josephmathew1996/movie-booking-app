package models

type PageDetail struct {
	PageNo   int `json:"page_number"`
	PageSize int `json:"page_size"`
}

type ErrorResponse struct {

	// Code describes the http status code
	// example: 400
	// in: number
	Code int `json:"code,omitempty"`

	// Message that describes the error occurred on server side
	// in: string
	Message string `json:"message,omitempty"`
}
