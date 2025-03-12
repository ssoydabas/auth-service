package dto

type StandardResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorData  `json:"error,omitempty"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"currentPage"`
	PageSize    int         `json:"pageSize"`
	TotalItems  int64       `json:"totalItems"`
	TotalPages  int         `json:"totalPages"`
}

type AccountResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PhotoUrl  string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
