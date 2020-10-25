package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/task_checker_api/src/controller"
	"github.com/task_checker_api/src/model"
	"github.com/task_checker_api/src/utils"
)

// SignUpHandler サインアップハンドラー
var SignUpHandler = func(w http.ResponseWriter, r *http.Request) {
	user, errorObj := SignUp(w, r)
	if errorObj.Message != "" {
		utils.Respond(w, http.StatusBadRequest, errorObj)
	} else {
		token := GetToken(user)
		// JWTを返却
		w.Write([]byte(token))
	}
}

// SignUp 新規登録処理
var SignUp = func(w http.ResponseWriter, r *http.Request) (model.Users, model.Error) {

	var signUp model.SignUp
	var errorObj model.Error
	var users model.Users

	json.NewDecoder(r.Body).Decode(&signUp)

	if signUp.MailAddress == "" {
		errorObj.Message = "\"MailAddress\" is missing"
		return users, errorObj
	}
	if signUp.UserName == "" {
		errorObj.Message = "\"UserName\" is missing"
		return users, errorObj
	}
	if signUp.Password == "" {
		errorObj.Message = "\"password\" is missing"
		return users, errorObj
	}

	isExist := checkExistMailAddress(signUp.MailAddress)

	if isExist {
		errorObj.Message = "isExist mailAddress"
	} else {
		CreateUser(signUp)
		users = controller.FindUserByMail(signUp.MailAddress)
	}
	return users, errorObj
}

// メールアドレス重複チェック
func checkExistMailAddress(mailAddress string) bool {
	db := model.DBConnect()
	result, err := db.Query("SELECT mail_address FROM users WHERE mail_address = $1;", mailAddress)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		var mail string
		err = result.Scan(&mail)
		if err != nil {
			panic(err.Error())
		}
		if mail != "" {
			return true
		}
	}
	return false
}

// CreateUser 登録
func CreateUser(signUp model.SignUp) {
	db := model.DBConnect()
	now := time.Now()

	hashPass := utils.UserPassHash(signUp.MailAddress, signUp.Password)

	_, err := db.Exec(`insert into users (
		user_name, password, mail_address, created_id, created_at, updated_id, updated_at) values ($1,$2,$3,$4,$5,$6,$7);`,
		signUp.UserName, hashPass, signUp.MailAddress, signUp.MailAddress, now, signUp.MailAddress, now)

	if err != nil {
		panic(err.Error())
	}
}
