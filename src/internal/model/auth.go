package model

type SignUp struct {
	Name       string `json:"name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokens struct {
	ID       int    `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}
