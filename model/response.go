package model

// Response represents the common API response structure
type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Data         any    `json:"data"`
	Total        int    `json:"total"`
}
