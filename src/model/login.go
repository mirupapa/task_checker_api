package model

//Login login
type Login struct {
	MailAddress string `form:"mailAddress" json:"mailAddress" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
}
