package main

import (
	"log"
	"net/http"
	"os"

	"main/src/auth"
	"main/src/controller"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	// ENVLoad()
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
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.GetTasksHandler),
		)).Methods("GET")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PostTask),
		)).Methods("POST")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PutTask),
		)).Methods("PUT")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.DeleteTask),
		)).Methods("DELETE")
		task.Handle("/done", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PutDone),
		)).Methods("PUT")
		task.Handle("/upSort", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.UpSort),
		)).Methods("PUT")
	}

	// cors
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "PUT", "DELETE"})

	//サーバー起動
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
