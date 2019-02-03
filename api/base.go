// Package api provides all useful API for frontend.
// A frontend should call required API through HTTP requests.
package api

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
