package models

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenExistsResponse struct {
	Exists       bool
	RefreshToken string
}
