package controller

import (
	"time"

	"main/src/model"
	"main/src/utils"
)

//FindUserByMail User検索
func FindUserByMail(mail string) model.Users {
	db := model.DBConnect()
	result, err := db.Query("SELECT mail_address, user_name FROM users WHERE mail_address = $1;", mail)
	if err != nil {
		panic(err.Error())
	}
	user := model.Users{}
	for result.Next() {
		var mailAddress, userName string
		err = result.Scan(&mailAddress, &userName)
		if err != nil {
			panic(err.Error())
		}
		user.MailAddress = mailAddress
		user.UserName = userName
	}
	return user
}

//UserCreate 登録 つくりかけ
func UserCreate(userMail string, users model.Users) {
	db := model.DBConnect()
	now := time.Now()

	hashPass := utils.UserPassHash(users.MailAddress, users.Password)

	_, err := db.Exec("INSERT INTO users (user_name, password, mail_address, created_id, created_at, updated_id, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)",
		users.UserName, hashPass, users.MailAddress, userMail, now, userMail, now)
	if err != nil {
		panic(err.Error())
	}
}
