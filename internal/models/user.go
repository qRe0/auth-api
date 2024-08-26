package models

type User struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	Password string  `json:"password,omitempty"`
	Rating   float64 `json:"rating"`
}
