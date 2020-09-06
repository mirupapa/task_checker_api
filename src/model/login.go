package model

//Login login
type Login struct {
	UserID   string `form:"userID" json:"userID" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
