package model

type Login struct {
	UserId   string `form:"userId" json:"userId" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
