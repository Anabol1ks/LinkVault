package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserRegisterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
