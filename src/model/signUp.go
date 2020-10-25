package model

//Login login
type SignUp struct {
	MailAddress string `form:"mailAddress" json:"mailAddress" binding:"required"`
	UserName    string `form:"userName" json:"userName" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}
