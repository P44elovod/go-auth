package models

type User struct {
	ID       int    `json:"-" pg:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" pg:"password_hash"`
}
