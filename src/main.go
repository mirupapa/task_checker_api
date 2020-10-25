package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/task_checker_api/src/auth"
	"github.com/task_checker_api/src/controller"
	"github.com/urfave/negroni"
)

// ENVLoad Env load
func ENVLoad() {
	err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		print("error_env")
	}
}

func main() {
	ENVLoad()
	r := mux.NewRouter()
	// ログイン
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	// サインアップ
	r.HandleFunc("/signUp", auth.SignUpHandler).Methods("POST")
	// 承認
	r.Handle("/auth", negroni.New(
		negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(auth.ExportUserInfo),
	))

	task := r.PathPrefix("/task").Subrouter()
	{
		task.Handle("/", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.GetTasks),
		))
	}

	// cors
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET"})

	//サーバー起動
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
