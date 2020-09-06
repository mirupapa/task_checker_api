package controller

import (
	"github.com/task_checker_api/src/model"
)

//LoginCheck ログイン判定
func LoginCheck(login model.Login) model.Users {
	db := model.DBConnect()
	result, err := db.Query("SELECT user_id, user_name FROM users WHERE user_id = $1 AND password = $2;", login.UserID, login.Password)
	//result, err := db.Query("select * from users;")

	if err != nil {
		panic(err.Error())
	}
	user := model.Users{}
	for result.Next() {
		var userID string
		var userName string
		err = result.Scan(&userID, &userName)
		if err != nil {
			panic(err.Error())
		}
		user.UserID = userID
		user.UserName = userName
	}
	return user
}
