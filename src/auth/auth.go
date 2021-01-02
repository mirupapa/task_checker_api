package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"main/src/controller"
	"main/src/model"
	"main/src/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	// dotenv
)

// LoginHandler ログインハンドラー
var LoginHandler = func(w http.ResponseWriter, r *http.Request) {
	user, errorObj := Login(w, r)
	if errorObj.Message != "" {
		utils.Respond(w, http.StatusBadRequest, errorObj)
	} else {
		token := GetToken(user)
		user.Token = token
		// JWTを返却
		// w.Write([]byte(result))
		json.NewEncoder(w).Encode(user)
	}
}

// Login ログイン処理
var Login = func(w http.ResponseWriter, r *http.Request) (model.Users, model.Error) {

	var login model.Login
	var errorObj model.Error
	var user model.Users

	json.NewDecoder(r.Body).Decode(&login)

	if login.MailAddress == "" {
		errorObj.Message = "\"MailAddress\" is missing"
		return user, errorObj
	}
	if login.Password == "" {
		errorObj.Message = "\"password\" is missing"
		return user, errorObj
	}

	user = LoginCheck(login)

	if user.MailAddress == "" {
		errorObj.Message = "login error"
	}
	return user, errorObj
}

// GetToken get token
var GetToken = func(user model.Users) string {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = user.MailAddress
	claims["name"] = user.UserName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// 電子署名
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	return tokenString
}

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNINGKEY")), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

//LoginCheck ログイン判定
func LoginCheck(login model.Login) model.Users {
	hashPass := utils.UserPassHash(login.MailAddress, login.Password)
	fmt.Printf(hashPass)
	db := model.DBConnect()
	result, err := db.Query("SELECT mail_address, user_name, password FROM users WHERE mail_address = $1;", login.MailAddress)
	if err != nil {
		panic(err.Error())
	}
	user := model.Users{}
	for result.Next() {
		var mailAddress, userName, password string
		err = result.Scan(&mailAddress, &userName, &password)
		if err != nil {
			panic(err.Error())
		}
		err := bcrypt.CompareHashAndPassword([]byte(password), []byte(login.MailAddress+" "+login.Password))
		if err != nil {
			panic(err.Error())
		}
		user.MailAddress = mailAddress
		user.UserName = userName
	}
	return user
}

// ExportUserInfo 承認後ユーザー情報返却
var ExportUserInfo = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	user := controller.FindUserByMail(mailAddress)
	// w.Write([]byte(mailAddress))
	json.NewEncoder(w).Encode(user)
})
