package request

import validation "github.com/go-ozzo/ozzo-validation"

// CreateCarReq represent create car request body
type CreateCarReq struct {
	Make           string `json:"make"`
	Model          string `json:"model"`
	Package        string `json:"package"`
	Color          string `json:"color"`
	Year           int    `json:"year"`
	Category       string `json:"category"`
	Mileage        int    `json:"mileage"`
	Price          int    `json:"price"`
	Identification string `json:"identification"`
}

func (request CreateCarReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Make, validation.Required),
		validation.Field(&request.Model, validation.Required),
		validation.Field(&request.Package, validation.Required),
		validation.Field(&request.Color, validation.Required),
		validation.Field(&request.Year, validation.Required),
		validation.Field(&request.Category, validation.Required),
		validation.Field(&request.Mileage, validation.Required),
		validation.Field(&request.Price, validation.Required),
		validation.Field(&request.Identification, validation.Required),
	)
}

// UpdateCarReq represent update car request body
type UpdateCarReq CreateCarReq

func (request UpdateCarReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Make, validation.Required),
		validation.Field(&request.Model, validation.Required),
		validation.Field(&request.Package, validation.Required),
		validation.Field(&request.Color, validation.Required),
		validation.Field(&request.Year, validation.Required),
		validation.Field(&request.Category, validation.Required),
		validation.Field(&request.Mileage, validation.Required),
		validation.Field(&request.Price, validation.Required),
		validation.Field(&request.Identification, validation.Required))
}
