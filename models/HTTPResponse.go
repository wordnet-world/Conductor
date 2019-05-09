package models

import (
	"encoding/json"
)

//This

// HTTPResponse is a standard response to a query from the
// DataService
type HTTPResponse struct {
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
	Success bool        `json:"success"`
	Count   int         `json:"count"`
}

// ToJSON returns the Json string representation
func (response HTTPResponse) ToJSON() string {
	res, _ := json.Marshal(response)
	return string(res)
}

// CreateHTTPResponse creates the response from an error passed in
func CreateHTTPResponse(error interface{}, data interface{}, success bool) HTTPResponse {

	return HTTPResponse{
		Data:    data,
		Error:   error,
		Success: success,
		Count:   0,
	}
}
