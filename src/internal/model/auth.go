package model

type SignUp struct {
	Name       string `json:"name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}
