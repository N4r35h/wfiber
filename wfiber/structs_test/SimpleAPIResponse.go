package structs_test

type SimpleAPIResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SimpleAPIResponseWithData[T any] struct {
	SimpleAPIResponse
	Data T `json:"data"`
}
