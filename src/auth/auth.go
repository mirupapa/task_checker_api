package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/task_checker_api/src/controller"
	"github.com/task_checker_api/src/model"
	"github.com/task_checker_api/src/utils"

	// dotenv
	"github.com/joho/godotenv"
)

// LoginHandler handler
var LoginHandler = func(w http.ResponseWriter, r *http.Request) {
	user, errorObj := Login(w, r)
	if errorObj.Message != "" {
		utils.Respond(w, http.StatusBadRequest, errorObj)
	} else {
		token := GetToken(user)
		// JWTを返却
		w.Write([]byte(token))
	}
}

// Login ログイン処理
var Login = func(w http.ResponseWriter, r *http.Request) (model.Users, model.Error) {

	var login model.Login
	var errorObj model.Error
	var user model.Users

	json.NewDecoder(r.Body).Decode(&login)

	if login.UserID == "" {
		errorObj.Message = "\"UserID\" is missing"
		return user, errorObj
	}
	if login.Password == "" {
		errorObj.Message = "\"password\" is missing"
		return user, errorObj
	}

	user = controller.LoginCheck(login)

	if user.UserID == "" {
		errorObj.Message = "login error"
	}
	return user, errorObj
}

// GetToken get token
var GetToken = func(user model.Users) string {
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		print("error_env")
	}
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = user.UserID
	claims["name"] = user.UserName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	print(os.Getenv("SIGNINGKEY"))

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
